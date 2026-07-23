package part

import (
	"context"

	"inventory/internal/model"
)

// ListParts возвращает список деталей с фильтрацией
func (s *PartService) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	// Получаем все детали
	allParts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Если фильтр пустой, возвращаем все
	if filter == nil {
		return allParts, nil
	}

	// Применяем фильтры
	return s.applyFilters(allParts, filter), nil
}

// applyFilters применяет фильтры к списку деталей
func (s *PartService) applyFilters(parts []*model.Part, filter *model.PartsFilter) []*model.Part {
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
