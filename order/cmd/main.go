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

	order_v1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*order_v1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*order_v1.OrderDto),
	}
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

// OrdersOrderUUIDCancelPost implements POST /orders/{order_uuid}/cancel operation.
// Отменяет заказ, если он ещё не оплачен.
func (o *OrderHandler) OrdersOrderUUIDCancelPost(ctx context.Context,
	params order_v1.OrdersOrderUUIDCancelPostParams,
) (order_v1.OrdersOrderUUIDCancelPostRes, error) {
	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	switch ord.Status {
	case order_v1.OrderStatusPAID:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot cancel already paid order",
		}, nil
	case order_v1.OrderStatusCANCELLED:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Order already cancelled",
		}, nil
	case order_v1.OrderStatusFULFILLED:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot cancel fulfilled order",
		}, nil
	}

	// Отменяем заказ
	if ord.Status == order_v1.OrderStatusPENDINGPAYMENT {
		ord.Status = order_v1.OrderStatusCANCELLED
	}

	return &order_v1.OrdersOrderUUIDCancelPostNoContent{}, nil
}

// OrdersOrderUUIDGet  implements GET /orders/{order_uuid} operation.
// Возвращает информацию о заказе по его UUID.
func (o *OrderHandler) OrdersOrderUUIDGet(ctx context.Context,
	params order_v1.OrdersOrderUUIDGetParams,
) (order_v1.OrdersOrderUUIDGetRes, error) {
	o.storage.mu.RLock()
	defer o.storage.mu.RUnlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	return ord, nil
}

// OrdersOrderUUIDPayPost implements POST /orders/{order_uuid}/pay operation.
// Проводит оплату ранее созданного заказа.
func (o *OrderHandler) OrdersOrderUUIDPayPost(ctx context.Context, req *order_v1.PayOrderRequest,
	params order_v1.OrdersOrderUUIDPayPostParams,
) (order_v1.OrdersOrderUUIDPayPostRes, error) {
	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()

	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid order UUID format",
		}, nil
	}

	ord, exists := o.storage.orders[orderUUID.String()]
	if !exists {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID),
		}, nil
	}

	switch ord.Status {
	case order_v1.OrderStatusPAID:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Order already paid",
		}, nil
	case order_v1.OrderStatusCANCELLED:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot pay cancelled order",
		}, nil
	case order_v1.OrderStatusFULFILLED:
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Cannot pay fulfilled order",
		}, nil
	}

	if req.PaymentMethod != "CARD" {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Payment method not card",
		}, nil
	}

	// Добавить интеграцию: вызывает PaymentService.PayOrder, передаёт user_uuid, order_uuid и payment_method. Получаетtransaction_uuid.

	// Обновляем статус заказа
	ord.Status = order_v1.OrderStatusPAID
	ord.PaymentMethod = order_v1.OptPaymentMethod{
		Value: req.PaymentMethod,
		Set:   true,
	}
	ord.TransactionUUID = order_v1.OptNilUUID{
		Value: uuid.New(), // Генерируем UUID транзакции
		Set:   true,
		Null:  false,
	}

	return &order_v1.PayOrderResponse{
		TransactionUUID: ord.TransactionUUID.Value,
	}, nil
}

// OrdersPost implements POST /orders operation.
// Создаёт новый заказ на основе выбранных пользователем деталей.
func (o *OrderHandler) OrdersPost(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.OrdersPostRes, error) {
	if len(req.PartUuids) == 0 {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "At least one part must be specified",
		}, nil
	}

	//  Добавить интеграцию
	//  Получает детали через InventoryService.ListParts.
	//	Проверяет, что все детали существуют. Если хотя бы одной нет — возвращает ошибку.
	//	Считает total_price.
	//	Генерирует order_uuid.
	//	Сохраняет заказ со статусом PENDING_PAYMENT.
	newOrder := &order_v1.OrderDto{
		OrderUUID:       uuid.New(),
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0, // добавить расчет по инфе из инвентори
		TransactionUUID: order_v1.OptNilUUID{Set: false},
		PaymentMethod:   order_v1.OptPaymentMethod{Set: false},
		Status:          order_v1.OrderStatusPENDINGPAYMENT,
	}

	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()
	o.storage.orders[newOrder.OrderUUID.String()] = newOrder

	return &order_v1.CreateOrderResponse{
		OrderUUID:  newOrder.OrderUUID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}

func main() {
	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage)
	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
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
		Addr:              net.JoinHostPort("0.0.0.0", httpPort), // <- Используем 0.0.0.0 вместо localhost
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s\n", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("Сервер остановлен")
}
