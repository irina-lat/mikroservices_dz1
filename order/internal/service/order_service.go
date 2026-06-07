package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"order/internal/client"
	"order/internal/repository"
	"order/pkg/model"
)

var (
	ErrOrderAlreadyPaid     = errors.New("order already paid")
	ErrOrderAlreadyCanceled = errors.New("order already canceled")
	ErrPartNotFound         = errors.New("some parts not found")
)

type OrderService struct {
	repo            repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

func NewOrderService(
	repo repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
) *OrderService {
	return &OrderService{
		repo:            repo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	// 1. Получаем детали через InventoryService.ListParts
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}

	// 2. Проверяем, что все детали существуют
	if len(parts) != len(partUUIDs) {
		return nil, ErrPartNotFound
	}

	// 3. Считаем total_price
	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	// 4. Генерируем order_uuid
	orderUUID := uuid.New().String()

	// 5. Сохраняем заказ со статусом PENDING_PAYMENT
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err := s.repo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// GetOrder возвращает заказ по UUID
func (s *OrderService) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.repo.FindByUUID(orderUUID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// PayOrder оплачивает заказ
func (s *OrderService) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	// 1. Находим заказ
	order, err := s.repo.FindByUUID(orderUUID)
	if err != nil {
		return "", err
	}

	// 2. Проверяем статус
	if order.Status == model.StatusPaid {
		return "", ErrOrderAlreadyPaid
	}
	if order.Status == model.StatusCancelled {
		return "", ErrOrderAlreadyCanceled
	}

	// 3. Вызываем PaymentService.PayOrder
	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		return "", fmt.Errorf("payment failed: %w", err)
	}

	// 4. Обновляем заказ
	method := model.PaymentMethod(paymentMethod)
	order.Status = model.StatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	if err := s.repo.Update(order); err != nil {
		return "", fmt.Errorf("failed to update order: %w", err)
	}

	return transactionUUID, nil
}

// CancelOrder отменяет заказ
func (s *OrderService) CancelOrder(ctx context.Context, orderUUID string) error {
	// 1. Находим заказ
	order, err := s.repo.FindByUUID(orderUUID)
	if err != nil {
		return err
	}

	// 2. Проверяем статус
	if order.Status == model.StatusPaid {
		return ErrOrderAlreadyPaid
	}

	// 3. Меняем статус на CANCELLED
	order.Status = model.StatusCancelled

	return s.repo.Update(order)
}