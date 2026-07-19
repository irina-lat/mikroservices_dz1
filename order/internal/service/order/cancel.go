package order

import (
	"context"
	"fmt"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// CancelOrder отменяет заказ
func (s *OrderService) CancelOrder(ctx context.Context, orderUUID string) error {
	// 1. Находим заказ
	order, err := s.repo.FindByUUID(ctx, orderUUID)
	if err != nil {
		return err
	}

	// 2. Проверяем статус
	if order.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	// 3. Меняем статус на CANCELLED
	order.Status = model.StatusCancelled

	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}
