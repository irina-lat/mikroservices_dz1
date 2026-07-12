package order

import (
	"context"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// FindByUUID находит заказ по UUID
func (r *InMemoryRepository) FindByUUID(ctx context.Context, uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return nil, model.ErrOrderNotFound
	}
	return order, nil
}
