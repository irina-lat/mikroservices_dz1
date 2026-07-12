package part

import (
	"context"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// FindByUUID находит деталь по UUID
func (r *InMemoryRepository) FindByUUID(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, model.ErrPartNotFound
	}
	return part, nil
}
