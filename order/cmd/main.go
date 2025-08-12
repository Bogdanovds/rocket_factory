package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

type OrderHandler struct {
	PaymentService   paymentV1.PaymentServiceClient
	InventoryService inventoryV1.InventoryServiceClient
	storage          *OrderStorage
}

func NewOrderHandler(storage *OrderStorage, paymentConn, inventoryConn *grpc.ClientConn) *OrderHandler {
	return &OrderHandler{
		storage:          storage,
		PaymentService:   paymentV1.NewPaymentServiceClient(paymentConn),
		InventoryService: inventoryV1.NewInventoryServiceClient(inventoryConn),
	}
}

// CancelOrder implements POST /orders/{order_uuid}/cancel operation.
// Отменяет заказ, если он ещё не оплачен.
func (o *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	switch ord.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot cancel already paid order",
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Order already cancelled",
		}, nil
	case orderV1.OrderStatusFULFILLED:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot cancel fulfilled order",
		}, nil
	}

	// Отменяем заказ
	ord.Status = orderV1.OrderStatusCANCELLED

	return &orderV1.CancelOrderNoContent{}, nil
}

// GetOrder  implements GET /orders/{order_uuid} operation.
// Возвращает информацию о заказе по его UUID.
func (o *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	o.storage.mu.RLock()
	defer o.storage.mu.RUnlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	return ord, nil
}

// PayOrder implements POST /orders/{order_uuid}/pay operation.
// Проводит оплату ранее созданного заказа.
func (o *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	switch ord.Status {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Order already paid",
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot pay cancelled order",
		}, nil
	case orderV1.OrderStatusFULFILLED:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot pay fulfilled order",
		}, nil
	}

	if req.PaymentMethod == "" {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Payment method is required",
		}, nil
	}

	resp, err := o.PaymentService.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      ord.UserUUID.String(),
		PaymentMethod: convertPaymentMethodToProto(req.PaymentMethod),
	})
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	// Обновляем статус заказа
	ord.Status = orderV1.OrderStatusPAID
	ord.PaymentMethod = orderV1.OptPaymentMethod{
		Value: req.PaymentMethod,
		Set:   true,
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction UUID format: %w", err)
	}

	ord.TransactionUUID = orderV1.OptNilUUID{
		Value: transactionUUID,
		Set:   true,
		Null:  false,
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: ord.TransactionUUID.Value,
	}, nil
}

func convertPaymentMethodToProto(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// CreateOrder implements POST /orders operation.
// Создаёт новый заказ на основе выбранных пользователем деталей.
func (o *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "At least one part must be specified",
		}, nil
	}

	partUUIDStrings := make([]string, len(req.PartUuids))
	for i, uuidVal := range req.PartUuids {
		partUUIDStrings[i] = uuidVal.String()
	}
	partsResp, err := o.InventoryService.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUUIDStrings,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get parts from inventory: %w", err)
	}

	if len(partsResp.Parts) != len(req.PartUuids) {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Some parts not found in inventory",
		}, nil
	}

	totalPrice := float32(0)
	for _, part := range partsResp.Parts {
		totalPrice += float32(part.Price)
	}

	newOrder := &orderV1.OrderDto{
		OrderUUID:       uuid.New(),
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: orderV1.OptNilUUID{Set: false},
		PaymentMethod:   orderV1.OptPaymentMethod{Set: false},
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
	}

	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()
	o.storage.orders[newOrder.OrderUUID.String()] = newOrder

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}

func main() {
	// PaymentService
	paymentConn, err := grpc.NewClient(
		"localhost:"+paymentServicePort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment service: %v", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v", err)
		}
	}()

	// InventoryService
	inventoryConn, err := grpc.NewClient(
		"localhost:"+inventoryServicePort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %v", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v", err)
		}
	}()

	// OrderService
	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage, paymentConn, inventoryConn)
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("failed to create OpenAPI server: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/v1")
				next.ServeHTTP(w, r)
			})
		}).Mount("/", orderServer)
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("HTTP server started on port %s\n", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v\n", err)
	}

	log.Println("Server stopped")
}
