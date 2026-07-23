package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *OrderServiceTestSuite) TestCancelOrder_Success() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()

	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  uuid.New().String(),
		PartUUIDs: []string{uuid.New().String()},
		Status:    model.StatusPendingPayment,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)
	s.repo.On("Update", ctx, order).Return(nil)

	// Вызов
	err := s.service.CancelOrder(ctx, orderUUID)

	// Проверки
	s.NoError(err)
	s.Equal(model.StatusCancelled, order.Status)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyPaid() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()

	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  uuid.New().String(),
		PartUUIDs: []string{uuid.New().String()},
		Status:    model.StatusPaid,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)

	// Вызов
	err := s.service.CancelOrder(ctx, orderUUID)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderAlreadyPaid, err)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCancelOrder_NotFound() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()

	s.repo.On("FindByUUID", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Вызов
	err := s.service.CancelOrder(ctx, orderUUID)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderNotFound, err)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCancelOrder_RepositoryError() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	expectedErr := errors.New("repository error")

	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  uuid.New().String(),
		PartUUIDs: []string{uuid.New().String()},
		Status:    model.StatusPendingPayment,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)
	s.repo.On("Update", ctx, order).Return(expectedErr)

	// Вызов
	err := s.service.CancelOrder(ctx, orderUUID)

	// Проверки
	s.Error(err)
	s.Contains(err.Error(), "failed to cancel order")

	s.repo.AssertExpectations(s.T())
}
