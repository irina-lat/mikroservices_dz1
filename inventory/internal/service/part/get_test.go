package part

import (
	"context"

	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

func (s *PartServiceTestSuite) TestGetPart_Success() {
	// Подготовка
	ctx := context.Background()
	partUUID := "test-uuid-123"
	expectedPart := &model.Part{
		UUID:          partUUID,
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.0,
		StockQuantity: 10,
		Category:      model.CategoryEngine,
	}

	s.repo.On("FindByUUID", ctx, partUUID).Return(expectedPart, nil)

	// Вызов
	part, err := s.service.GetPart(ctx, partUUID)

	// Проверки
	s.NoError(err)
	s.Equal(expectedPart, part)
	s.Equal(partUUID, part.UUID)
	s.Equal("Test Part", part.Name)
	s.Equal(100.0, part.Price)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestGetPart_NotFound() {
	// Подготовка
	ctx := context.Background()
	partUUID := "non-existent-uuid"

	s.repo.On("FindByUUID", ctx, partUUID).Return(nil, model.ErrPartNotFound)

	// Вызов
	part, err := s.service.GetPart(ctx, partUUID)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrPartNotFound, err)
	s.Nil(part)

	s.repo.AssertExpectations(s.T())
}

func (s *PartServiceTestSuite) TestGetPart_EmptyUUID() {
	// Подготовка
	ctx := context.Background()
	partUUID := ""

	s.repo.On("FindByUUID", ctx, partUUID).Return(nil, model.ErrPartNotFound)

	// Вызов
	part, err := s.service.GetPart(ctx, partUUID)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrPartNotFound, err)
	s.Nil(part)

	s.repo.AssertExpectations(s.T())
}
