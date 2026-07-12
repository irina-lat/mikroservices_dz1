package order

import (
	"context"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// Save сохраняет заказ в памяти
func (r *InMemoryRepository) Save(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.OrderUUID] = order
	return nil
}
