package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/irina-lat/microservices-course/inventory/internal/api/inventory/v1"
	"github.com/irina-lat/microservices-course/inventory/internal/repository/part"
	partsvc "github.com/irina-lat/microservices-course/inventory/internal/service/part"

	pb "shared/pkg/proto/inventory/v1"
)

func main() {
	// 1. Создаём репозиторий
	repo := part.NewInMemoryRepository()

	// 2. Инициализируем тестовые данные
	ctx := context.Background()
	if err := repo.InitSampleData(ctx); err != nil {
		log.Printf("Warning: failed to init sample data: %v", err)
	}

	// 3. Создаём сервис
	svc := partsvc.NewService(repo)

	// 4. Создаём gRPC API
	api := apiv1.NewAPI(svc)

	// 5. Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterInventoryServiceServer(grpcServer, api)

	// 6. Запускаем сервер
	log.Println("InventoryService starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}