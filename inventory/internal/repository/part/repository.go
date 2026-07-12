package part

import (
	"context"
	"sync"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// Repository определяет интерфейс для работы с деталями
type Repository interface {
	Save(ctx context.Context, part *model.Part) error
	FindByUUID(ctx context.Context, uuid string) (*model.Part, error)
	FindAll(ctx context.Context) ([]*model.Part, error)
}

// InMemoryRepository реализует Repository в памяти
type InMemoryRepository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

// NewInMemoryRepository создаёт новый экземпляр репозитория
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		parts: make(map[string]*model.Part),
	}
}

// Save сохраняет деталь в репозитории
func (r *InMemoryRepository) Save(ctx context.Context, part *model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.parts[part.UUID] = part
	return nil
}
