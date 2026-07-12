package order

import (
	"context"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// Update обновляет заказ в памяти
func (r *InMemoryRepository) Update(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return model.ErrOrderNotFound
	}
	r.orders[order.OrderUUID] = order
	return nil
}
