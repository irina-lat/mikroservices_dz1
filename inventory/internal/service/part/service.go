package part

import (
	"context"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// Service определяет интерфейс бизнес-логики для работы с деталями
type Service interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}

// Repository определяет интерфейс для работы с хранилищем
type Repository interface {
	FindByUUID(ctx context.Context, uuid string) (*model.Part, error)
	FindAll(ctx context.Context) ([]*model.Part, error)
}

// PartService реализует бизнес-логику
type PartService struct {
	repo Repository
}

// NewService создаёт новый экземпляр PartService
func NewService(repo Repository) *PartService {
	return &PartService{
		repo: repo,
	}
}
