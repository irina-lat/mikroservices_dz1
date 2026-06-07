package main

import (
	"log"
	"net"

	"payment/internal/handler"
	"payment/internal/service"
	pb "shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Создаём сервис
	paymentService := service.NewPaymentService()

	// 2. Создаём gRPC хендлер
	grpcHandler := handler.NewGrpcHandler(paymentService)

	// 3. Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterPaymentServiceServer(grpcServer, grpcHandler)

	// 4. Запускаем сервер
	log.Println("PaymentService starting on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}