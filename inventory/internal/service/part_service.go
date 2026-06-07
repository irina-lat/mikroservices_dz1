package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"inventory/internal/repository"
	"inventory/pkg/model"
)

type PartService struct {
	repo repository.PartRepository
}

func NewPartService(repo repository.PartRepository) *PartService {
	return &PartService{
		repo: repo,
	}
}

// GetPart возвращает деталь по UUID
func (s *PartService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	return s.repo.FindByUUID(uuid)
}

// ListParts возвращает список деталей с фильтрацией
func (s *PartService) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	// Получаем все детали
	allParts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Применяем фильтры
	filtered := s.applyFilters(allParts, filter)

	return filtered, nil
}

// applyFilters применяет фильтры к списку деталей
func (s *PartService) applyFilters(parts []*model.Part, filter *model.PartsFilter) []*model.Part {
	if filter == nil {
		return parts
	}

	result := make([]*model.Part, 0, len(parts))

	for _, part := range parts {
		if s.matchesFilter(part, filter) {
			result = append(result, part)
		}
	}

	return result
}

// matchesFilter проверяет, соответствует ли деталь фильтру
func (s *PartService) matchesFilter(part *model.Part, filter *model.PartsFilter) bool {
	// Фильтр по UUID (OR внутри поля)
	if len(filter.UUIDs) > 0 {
		found := false
		for _, uuid := range filter.UUIDs {
			if part.UUID == uuid {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по имени (OR внутри поля)
	if len(filter.Names) > 0 {
		found := false
		for _, name := range filter.Names {
			if part.Name == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по категории (OR внутри поля)
	if len(filter.Categories) > 0 {
		found := false
		for _, category := range filter.Categories {
			if part.Category == category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по стране производителя (OR внутри поля)
	if len(filter.ManufacturerCountries) > 0 {
		found := false
		for _, country := range filter.ManufacturerCountries {
			if part.Manufacturer.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по тегам (OR внутри поля)
	if len(filter.Tags) > 0 {
		found := false
		for _, tag := range filter.Tags {
			for _, partTag := range part.Tags {
				if partTag == tag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// CreateSampleParts создаёт тестовые детали
func (s *PartService) CreateSampleParts() error {
	now := time.Now()

	sampleParts := []*model.Part{
		{
			UUID:          uuid.New().String(),
			Name:          "Main Engine",
			Description:   "Main propulsion engine for spacecraft",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions: model.Dimensions{
				Length: 200.0,
				Width:  150.0,
				Height: 100.0,
				Weight: 500.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "https://spacex.com",
			},
			Tags:      []string{"engine", "propulsion", "main"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Fuel Tank",
			Description:   "High-capacity fuel tank",
			Price:         5000.0,
			StockQuantity: 10,
			Category:      model.CategoryFuel,
			Dimensions: model.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 200.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Boeing",
				Country: "USA",
				Website: "https://boeing.com",
			},
			Tags:      []string{"fuel", "tank", "storage"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Porthole Window",
			Description:   "Reinforced viewing window",
			Price:         2500.0,
			StockQuantity: 20,
			Category:      model.CategoryPorthole,
			Dimensions: model.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Thales",
				Country: "France",
				Website: "https://thalesgroup.com",
			},
			Tags:      []string{"window", "viewing", "porthole"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, part := range sampleParts {
		if err := s.repo.Save(part); err != nil {
			return err
		}
	}

	return nil
}