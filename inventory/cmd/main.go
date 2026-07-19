package main

import (
	"context"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/irina-lat/microservices-course/inventory/internal/api/inventory/v1"
	"github.com/irina-lat/microservices-course/inventory/internal/repository/part"
	partsvc "github.com/irina-lat/microservices-course/inventory/internal/service/part"

	pb "shared/pkg/proto/inventory/v1"
)

func main() {
	ctx := context.Background()

	// 1. Подключаемся к MongoDB
	uri := "mongodb://inventory-service-user:inventory-service-password@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("✅ Connected to MongoDB")

	db := client.Database("inventory-service")

	// 2. Создаём репозиторий (MongoDB)
	repo := part.NewMongoRepository(db) // ← изменено с NewInMemoryRepository на NewMongoRepository

	// 3. Инициализируем тестовые данные
	log.Println("📦 Initializing sample data...")
	if err := repo.InitSampleData(ctx); err != nil {
		log.Printf("Warning: failed to init sample data: %v", err)
	}

	// 4. Создаём сервис
	svc := partsvc.NewService(repo)

	// 5. Создаём gRPC API
	api := apiv1.NewAPI(svc)

	// 6. Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterInventoryServiceServer(grpcServer, api)

	// 7. Запускаем сервер
	log.Println("🚀 InventoryService starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}