package repository

import (
	"errors"
	"sync"

	"order/pkg/model"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository interface {
	Save(order *model.Order) error
	FindByUUID(uuid string) (*model.Order, error)
	Update(order *model.Order) error
}

type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

func (r *InMemoryOrderRepository) Save(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.OrderUUID] = order
	return nil
}

func (r *InMemoryOrderRepository) FindByUUID(uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (r *InMemoryOrderRepository) Update(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return ErrOrderNotFound
	}
	r.orders[order.OrderUUID] = order
	return nil
}