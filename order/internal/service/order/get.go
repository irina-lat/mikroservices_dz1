package order

import (
	"context"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// GetOrder возвращает заказ по UUID
func (s *OrderService) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	return s.repo.FindByUUID(ctx, orderUUID)
}
