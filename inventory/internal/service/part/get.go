package part

import (
	"context"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// GetPart возвращает деталь по UUID
func (s *PartService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	return s.repo.FindByUUID(ctx, uuid)
}
