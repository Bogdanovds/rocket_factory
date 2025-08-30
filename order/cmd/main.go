package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/bogdanovds/rocket_factory/order/internal/api/order/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/inventory/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/payment/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/repository/order"
	order2 "github.com/bogdanovds/rocket_factory/order/internal/service/order"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

const (
	httpPort             = "8080"
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
)

func main() {
	paymentConn := createGRPCConn(paymentServicePort)
	defer func(paymentConn *grpc.ClientConn) {
		err := paymentConn.Close()
		if err != nil {
			log.Printf("Error closing payment connection: %v", err)
		}
	}(paymentConn)

	inventoryConn := createGRPCConn(inventoryServicePort)
	defer func(inventoryConn *grpc.ClientConn) {
		err := inventoryConn.Close()
		if err != nil {
			log.Printf("Error closing inventory connection: %v", err)
		}
	}(inventoryConn)

	paymentClient := payment.New(paymentConn)
	inventoryClient := inventory.New(inventoryConn)

	// OrderService
	orderRepo := order.NewRepo()
	orderService := order2.NewService(orderRepo, inventoryClient, paymentClient)
	orderHandler := v1.NewHandler(orderService)

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

func createGRPCConn(port string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to gRPC service: %v", err)
	}
	return conn
}
