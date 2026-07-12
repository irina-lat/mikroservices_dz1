package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/irina-lat/microservices-course/payment/internal/api/payment/v1"
	"github.com/irina-lat/microservices-course/payment/internal/service/payment"

	pb "shared/pkg/proto/payment/v1"
)

func main() {
	// 1. Создаём сервис
	svc := payment.New()

	// 2. Создаём gRPC API
	api := apiv1.NewAPI(svc)

	// 3. Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterPaymentServiceServer(grpcServer, api)

	// 4. Запускаем сервер
	log.Println("PaymentService starting on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}