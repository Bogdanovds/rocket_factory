// package main
//
// import (
//
//	"context"
//	"fmt"
//	"log"
//	"net"
//	"os"
//	"os/signal"
//	"syscall"
//
//	"github.com/google/uuid"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/reflection"
//
//	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
//
// )
//
// const grpcPort = 50052
//
//	type PaymentService struct {
//		paymentV1.UnimplementedPaymentServiceServer
//	}
//
//	func NewPaymentService() *PaymentService {
//		return &PaymentService{}
//	}
//
// // PayOrder обрабатывает платеж и возвращает UUID транзакции
//
//	func (s *PaymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
//		transactionUUID := uuid.New().String()
//
//		log.Printf("Оплата прошла успешно, transaction_uuid: %s\n"+
//			"Детали платежа:\n"+
//			" - Order UUID: %s\n"+
//			" - User UUID: %s\n"+
//			" - Payment Method: %s",
//			transactionUUID, req.OrderUuid, req.UserUuid, req.PaymentMethod.String())
//
//		return &paymentV1.PayOrderResponse{
//			TransactionUuid: transactionUUID,
//		}, nil
//	}
//
//	func main() {
//		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
//		if err != nil {
//			log.Printf("failed to listen: %v\n", err)
//			return
//		}
//		defer func() {
//			if cerr := lis.Close(); cerr != nil {
//				log.Printf("failed to close listener: %v\n", cerr)
//			}
//		}()
//
//		s := grpc.NewServer()
//		paymentService := NewPaymentService()
//		paymentV1.RegisterPaymentServiceServer(s, paymentService)
//
//		reflection.Register(s)
//
//		go func() {
//			log.Printf("🚀 gRPC s listening on %d\n", grpcPort)
//			err = s.Serve(lis)
//			if err != nil {
//				log.Printf("failed to serve: %v\n", err)
//				return
//			}
//		}()
//
//		// Graceful shutdown
//		quit := make(chan os.Signal, 1)
//		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//		<-quit
//		log.Println("🛑 Shutting down gRPC s...")
//		s.GracefulStop()
//		log.Println("✅ s stopped")
//	}
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanovds/rocket_factory/payment/internal/service/payment"
	paymentV1 "github.com/bogdanovds/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

func main() {
	// Create payment service
	paymentService := payment.NewPaymentService()

	// Create gRPC server with API
	grpcServer := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(grpcServer, paymentService)
	reflection.Register(grpcServer)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Start server
	go func() {
		log.Printf("🚀 Payment gRPC server listening on port %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
