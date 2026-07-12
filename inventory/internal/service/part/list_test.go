package part

import (
	"context"
	"errors"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

func (s *PartServiceTestSuite) TestListParts_All() {
	// Подготовка
	ctx := context.Background()
	expectedParts := []*model.Part{
		{
			UUID:          "uuid-1",
			Name:          "Part 1",
			Price:         100.0,
			Category:      model.CategoryEngine,
			Manufacturer:  model.Manufacturer{Country: "USA"},
			Tags:          []string{"tag1", "tag2"},
		},
		{
			UUID:          "uuid-2",
			Name:          "Part 2",
			Price:         200.0,
			Category:      model.CategoryFuel,
			Manufacturer:  model.Manufacturer{Country: "Germany"},
			Tags:          []string{"tag3"},
		},
	}

	s.repo.On("FindAll", ctx).Return(expectedParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, nil)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByUUID() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		UUIDs: []string{"uuid-1"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Equal("uuid-1", parts[0].UUID)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleUUIDs() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		UUIDs: []string{"uuid-1", "uuid-2"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
		{UUID: "uuid-3", Name: "Part 3", Price: 300.0},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByName() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Names: []string{"Part 1"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Equal("Part 1", parts[0].Name)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleNames() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Names: []string{"Part 1", "Part 2"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
		{UUID: "uuid-3", Name: "Part 3", Price: 300.0},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Price: 100.0},
		{UUID: "uuid-2", Name: "Part 2", Price: 200.0},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByCategory() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Category: model.CategoryEngine},
		{UUID: "uuid-2", Name: "Part 2", Category: model.CategoryFuel},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Category: model.CategoryEngine},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Equal(model.CategoryEngine, parts[0].Category)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleCategories() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine, model.CategoryFuel},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Category: model.CategoryEngine},
		{UUID: "uuid-2", Name: "Part 2", Category: model.CategoryFuel},
		{UUID: "uuid-3", Name: "Part 3", Category: model.CategoryWing},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Category: model.CategoryEngine},
		{UUID: "uuid-2", Name: "Part 2", Category: model.CategoryFuel},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByManufacturerCountry() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		ManufacturerCountries: []string{"USA"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Manufacturer: model.Manufacturer{Country: "USA"}},
		{UUID: "uuid-2", Name: "Part 2", Manufacturer: model.Manufacturer{Country: "Germany"}},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Manufacturer: model.Manufacturer{Country: "USA"}},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Equal("USA", parts[0].Manufacturer.Country)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleCountries() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		ManufacturerCountries: []string{"USA", "Germany"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Manufacturer: model.Manufacturer{Country: "USA"}},
		{UUID: "uuid-2", Name: "Part 2", Manufacturer: model.Manufacturer{Country: "Germany"}},
		{UUID: "uuid-3", Name: "Part 3", Manufacturer: model.Manufacturer{Country: "France"}},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Manufacturer: model.Manufacturer{Country: "USA"}},
		{UUID: "uuid-2", Name: "Part 2", Manufacturer: model.Manufacturer{Country: "Germany"}},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByTags() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Tags: []string{"tag1"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Tags: []string{"tag1", "tag2"}},
		{UUID: "uuid-2", Name: "Part 2", Tags: []string{"tag3"}},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Tags: []string{"tag1", "tag2"}},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Contains(parts[0].Tags, "tag1")

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleTags() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Tags: []string{"tag1", "tag3"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Tags: []string{"tag1", "tag2"}},
		{UUID: "uuid-2", Name: "Part 2", Tags: []string{"tag3"}},
		{UUID: "uuid-3", Name: "Part 3", Tags: []string{"tag4"}},
	}
	expectedParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1", Tags: []string{"tag1", "tag2"}},
		{UUID: "uuid-2", Name: "Part 2", Tags: []string{"tag3"}},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterByMultipleFields() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		Categories:           []model.Category{model.CategoryEngine},
		ManufacturerCountries: []string{"USA"},
	}

	allParts := []*model.Part{
		{
			UUID:         "uuid-1",
			Name:         "Part 1",
			Category:     model.CategoryEngine,
			Manufacturer: model.Manufacturer{Country: "USA"},
		},
		{
			UUID:         "uuid-2",
			Name:         "Part 2",
			Category:     model.CategoryEngine,
			Manufacturer: model.Manufacturer{Country: "Germany"},
		},
		{
			UUID:         "uuid-3",
			Name:         "Part 3",
			Category:     model.CategoryFuel,
			Manufacturer: model.Manufacturer{Country: "USA"},
		},
	}
	expectedParts := []*model.Part{
		{
			UUID:         "uuid-1",
			Name:         "Part 1",
			Category:     model.CategoryEngine,
			Manufacturer: model.Manufacturer{Country: "USA"},
		},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Len(parts, 1)
	s.Equal(model.CategoryEngine, parts[0].Category)
	s.Equal("USA", parts[0].Manufacturer.Country)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_FilterNoMatch() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{
		UUIDs: []string{"non-existent-uuid"},
	}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1"},
		{UUID: "uuid-2", Name: "Part 2"},
	}
	expectedParts := []*model.Part{}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(expectedParts, parts)
	s.Empty(parts)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_EmptyFilter() {
	// Подготовка
	ctx := context.Background()
	filter := &model.PartsFilter{}

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1"},
		{UUID: "uuid-2", Name: "Part 2"},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, filter)

	// Проверки
	s.NoError(err)
	s.Equal(allParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_NilFilter() {
	// Подготовка
	ctx := context.Background()

	allParts := []*model.Part{
		{UUID: "uuid-1", Name: "Part 1"},
		{UUID: "uuid-2", Name: "Part 2"},
	}

	s.repo.On("FindAll", ctx).Return(allParts, nil)

	// Вызов
	parts, err := s.service.ListParts(ctx, nil)

	// Проверки
	s.NoError(err)
	s.Equal(allParts, parts)
	s.Len(parts, 2)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestListParts_RepositoryError() {
	// Подготовка
	ctx := context.Background()
	expectedErr := errors.New("repository error")

	s.repo.On("FindAll", ctx).Return(nil, expectedErr)

	// Вызов
	parts, err := s.service.ListParts(ctx, nil)

	// Проверки
	s.Error(err)
	s.Equal(expectedErr, err)
	s.Nil(parts)

	s.repo.AssertExpectations(s.T())
}
