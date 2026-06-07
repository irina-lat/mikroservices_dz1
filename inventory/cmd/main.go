package main

import (
	"log"
	"net"

	"inventory/internal/handler"
	"inventory/internal/repository"
	"inventory/internal/service"
	pb "shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Создаём репозиторий и сервис
	partRepo := repository.NewInMemoryPartRepository()
	partService := service.NewPartService(partRepo)

	// 2. Добавляем тестовые данные
	if err := partService.CreateSampleParts(); err != nil {
		log.Printf("Warning: failed to create sample parts: %v", err)
	}

	// 3. Создаём gRPC хендлер
	grpcHandler := handler.NewGrpcHandler(partService)

	// 4. Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterInventoryServiceServer(grpcServer, grpcHandler)

	// 5. Запускаем сервер
	log.Println("InventoryService starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}