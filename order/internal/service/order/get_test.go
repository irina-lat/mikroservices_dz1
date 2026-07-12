package order

import (
	"context"
	"errors"

	"github.com/irina-lat/microservices-course/order/internal/model"

	"github.com/google/uuid"
)

func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	expectedOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New().String(),
		PartUUIDs:  []string{uuid.New().String()},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(expectedOrder, nil)

	// Вызов
	order, err := s.service.GetOrder(ctx, orderUUID)

	// Проверки
	s.NoError(err)
	s.Equal(expectedOrder, order)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()

	s.repo.On("FindByUUID", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Вызов
	order, err := s.service.GetOrder(ctx, orderUUID)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderNotFound, err)
	s.Nil(order)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetOrder_RepositoryError() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	expectedErr := errors.New("repository error")

	s.repo.On("FindByUUID", ctx, orderUUID).Return(nil, expectedErr)

	// Вызов
	order, err := s.service.GetOrder(ctx, orderUUID)

	// Проверки
	s.Error(err)
	s.Equal(expectedErr, err)
	s.Nil(order)

	s.repo.AssertExpectations(s.T())
}
