package order

import (
	"context"

	"order/internal/model"
)

// GetOrder возвращает заказ по UUID
func (s *OrderService) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	return s.repo.FindByUUID(ctx, orderUUID)
}
