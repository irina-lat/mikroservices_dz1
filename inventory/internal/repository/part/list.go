package part

import (
	"context"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// FindAll возвращает все детали
func (r *InMemoryRepository) FindAll(ctx context.Context) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]*model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}
	return parts, nil
}
