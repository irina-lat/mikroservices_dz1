package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/irina-lat/microservices-course/order/internal/api/order/v1"
	inventoryv1 "github.com/irina-lat/microservices-course/order/internal/client/grpc/inventory/v1"
	paymentv1 "github.com/irina-lat/microservices-course/order/internal/client/grpc/payment/v1"
	"github.com/irina-lat/microservices-course/order/internal/migrator"
	"github.com/irina-lat/microservices-course/order/internal/repository/order"
	orderservice "github.com/irina-lat/microservices-course/order/internal/service/order"

	orderapi "shared/pkg/openapi/order/v1"
	inventorypb "shared/pkg/proto/inventory/v1"
	paymentpb "shared/pkg/proto/payment/v1"
)

func main() {
	ctx := context.Background()

	// 1. Подключаемся к PostgreSQL
	dsn := "host=localhost port=5432 user=order-service-user password=123 dbname=order-service sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// 2. Запускаем миграции
	db := stdlib.OpenDBFromPool(pool)
	if err := migrator.Run(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	db.Close()

	// 3. Подключаемся к InventoryService (gRPC)
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer inventoryConn.Close()
	log.Println("✅ Connected to InventoryService on :50051")

	// 4. Подключаемся к PaymentService (gRPC)
	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	defer paymentConn.Close()
	log.Println("✅ Connected to PaymentService on :50052")

	// 5. Создаём клиентов
	inventoryClient := inventoryv1.NewInventoryClient(inventorypb.NewInventoryServiceClient(inventoryConn))
	paymentClient := paymentv1.NewPaymentClient(paymentpb.NewPaymentServiceClient(paymentConn))

	// 6. Создаём репозиторий (PostgreSQL)
	orderRepo := order.NewPostgresRepository(pool)

	// 7. Создаём сервис
	orderService := orderservice.NewService(orderRepo, inventoryClient, paymentClient)

	// 8. Создаём API хендлеры
	orderHandler := apiv1.NewAPI(orderService)

	// 9. Создаём HTTP роутер
	router, err := orderapi.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	// 10. Запускаем HTTP сервер
	log.Println("🚀 OrderService starting on :8080")
	log.Println("Available endpoints:")
	log.Println("  POST   /api/v1/orders")
	log.Println("  GET    /api/v1/orders/{order_uuid}")
	log.Println("  POST   /api/v1/orders/{order_uuid}/pay")
	log.Println("  POST   /api/v1/orders/{order_uuid}/cancel")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}