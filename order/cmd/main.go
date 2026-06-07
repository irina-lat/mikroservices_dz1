package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"order/internal/client"
	"order/internal/handler"
	"order/internal/repository"
	"order/internal/service"
	orderapi "shared/pkg/openapi/order/v1"
	inventorypb "shared/pkg/proto/inventory/v1"
	paymentpb "shared/pkg/proto/payment/v1"
)

func main() {
	// 1. Подключаемся к InventoryService (gRPC)
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := client.NewGrpcInventoryClient(inventorypb.NewInventoryServiceClient(inventoryConn))

	// 2. Подключаемся к PaymentService (gRPC)
	paymentConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := client.NewGrpcPaymentClient(paymentpb.NewPaymentServiceClient(paymentConn))

	// 3. Создаём репозиторий, сервис и хендлер
	orderRepo := repository.NewInMemoryOrderRepository()
	orderService := service.NewOrderService(orderRepo, inventoryClient, paymentClient)
	orderHandler := handler.NewOrderHandler(orderService)

	// 4. Создаём HTTP сервер с Ogen роутером
	router, err := orderapi.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	// 5. Запускаем HTTP сервер
	log.Println("OrderService starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}