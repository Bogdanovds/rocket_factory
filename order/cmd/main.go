package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/bogdanovds/rocket_factory/order/internal/api/order/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/inventory/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/client/grpc/payment/v1"
	"github.com/bogdanovds/rocket_factory/order/internal/migrator"
	"github.com/bogdanovds/rocket_factory/order/internal/repository/postgres"
	order2 "github.com/bogdanovds/rocket_factory/order/internal/service/order"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

const (
	httpPort             = "8081"
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
)

func main() {
	// Подключение к PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Применение миграций
	m := migrator.New(db)
	if err := m.UpEmbed(); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

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

	// OrderService с PostgreSQL репозиторием
	orderRepo := postgres.NewRepository(db)
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

func connectDB() (*sql.DB, error) {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5433")
	user := getEnv("POSTGRES_USER", "order-service-user")
	password := getEnv("POSTGRES_PASSWORD", "order-service-password")
	dbname := getEnv("POSTGRES_DB", "order-service")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
