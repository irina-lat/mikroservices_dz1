package order

import (
	"context"
	"sync"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// Repository определяет интерфейс для работы с заказами
type Repository interface {
	Save(ctx context.Context, order *model.Order) error
	FindByUUID(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

// InMemoryRepository реализует Repository в памяти
type InMemoryRepository struct {
	orders map[string]*model.Order
	mu     sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		orders: make(map[string]*model.Order),
	}
}
