package order

import (
	"context"

	"order/internal/model"
)

// Service определяет интерфейс бизнес-логики заказов
type Service interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error)
	CancelOrder(ctx context.Context, orderUUID string) error
	UpdateOrderStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error // ← добавлен
}

// Repository определяет интерфейс для работы с хранилищем
type Repository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByUUID(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

// InventoryClient определяет интерфейс клиента InventoryService
type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]*model.Part, error)
}

// PaymentClient определяет интерфейс клиента PaymentService
type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (string, error)
}

// KafkaProducer определяет интерфейс для отправки событий в Kafka
type KafkaProducer interface {
	SendOrderPaid(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error
}

// OrderService реализует бизнес-логику
type OrderService struct {
	repo            Repository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	kafkaProducer   KafkaProducer // ← добавлен
}

// NewService создаёт новый экземпляр OrderService
func NewService(
	repo Repository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	kafkaProducer KafkaProducer, // ← добавлен
) *OrderService {
	return &OrderService{
		repo:            repo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		kafkaProducer:   kafkaProducer,
	}
}

// UpdateOrderStatus обновляет статус заказа
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error {
	order, err := s.repo.FindByUUID(ctx, orderUUID)
	if err != nil {
		return err
	}

	order.Status = status
	return s.repo.Update(ctx, order)
}