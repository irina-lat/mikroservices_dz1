package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/irina-lat/microservices-course/order/internal/api/order/v1"
	inventoryv1 "github.com/irina-lat/microservices-course/order/internal/client/grpc/inventory/v1"
	paymentv1 "github.com/irina-lat/microservices-course/order/internal/client/grpc/payment/v1"
	"github.com/irina-lat/microservices-course/order/internal/repository/order"
	orderservice "github.com/irina-lat/microservices-course/order/internal/service/order"

	orderapi "shared/pkg/openapi/order/v1"
	inventorypb "shared/pkg/proto/inventory/v1"
	paymentpb "shared/pkg/proto/payment/v1"
)

func main() {
	// 1. Подключаемся к InventoryService (gRPC, порт 50051)
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer inventoryConn.Close()
	log.Println("Connected to InventoryService on :50051")

	// 2. Подключаемся к PaymentService (gRPC, порт 50052)
	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	defer paymentConn.Close()
	log.Println("Connected to PaymentService on :50052")

	// 3. Создаём клиентов для gRPC сервисов
	inventoryClient := inventoryv1.NewInventoryClient(inventorypb.NewInventoryServiceClient(inventoryConn))
	paymentClient := paymentv1.NewPaymentClient(paymentpb.NewPaymentServiceClient(paymentConn))

	// 4. Создаём репозиторий (In-Memory)
	orderRepo := order.NewInMemoryRepository()

	// 5. Создаём сервис с бизнес-логикой
	orderService := orderservice.NewService(orderRepo, inventoryClient, paymentClient)

	// 6. Создаём API хендлеры
	orderHandler := apiv1.NewAPI(orderService)

	// 7. Создаём HTTP роутер из сгенерированного Ogen кода
	router, err := orderapi.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	// 8. Запускаем HTTP сервер
	log.Println("OrderService starting on :8080")
	log.Println("Available endpoints:")
	log.Println("  POST   /api/v1/orders")
	log.Println("  GET    /api/v1/orders/{order_uuid}")
	log.Println("  POST   /api/v1/orders/{order_uuid}/pay")
	log.Println("  POST   /api/v1/orders/{order_uuid}/cancel")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}