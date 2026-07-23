package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
)

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	// 1. Получаем детали через InventoryService
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}

	// 2. Проверяем, что все детали существуют
	if len(parts) != len(partUUIDs) {
		return nil, model.ErrPartNotFound
	}

	// 3. Считаем total_price
	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	// 4. Генерируем order_uuid
	orderUUID := uuid.New().String()

	// 5. Создаём заказ со статусом PENDING_PAYMENT
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	// 6. Сохраняем заказ
	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}
