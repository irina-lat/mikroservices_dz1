package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *OrderServiceTestSuite) TestPayOrder_Success() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	paymentMethod := "CARD"
	transactionUUID := uuid.New().String()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New().String(),
		PartUUIDs:  []string{uuid.New().String()},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", ctx, order.UserUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)
	s.repo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	// Вызов
	result, err := s.service.PayOrder(ctx, orderUUID, paymentMethod)

	// Проверки
	s.NoError(err)
	s.Equal(transactionUUID, result)
	s.Equal(model.StatusPaid, order.Status)
	s.Equal(&transactionUUID, order.TransactionUUID)
	expectedMethod := model.PaymentMethod(paymentMethod)
	s.Equal(&expectedMethod, order.PaymentMethod)

	s.repo.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_AlreadyPaid() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	paymentMethod := "CARD"

	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  uuid.New().String(),
		PartUUIDs: []string{uuid.New().String()},
		Status:    model.StatusPaid,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)

	// Вызов
	result, err := s.service.PayOrder(ctx, orderUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderAlreadyPaid, err)
	s.Empty(result)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_AlreadyCanceled() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	paymentMethod := "CARD"

	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  uuid.New().String(),
		PartUUIDs: []string{uuid.New().String()},
		Status:    model.StatusCancelled,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)

	// Вызов
	result, err := s.service.PayOrder(ctx, orderUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderAlreadyCanceled, err)
	s.Empty(result)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_NotFound() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	paymentMethod := "CARD"

	s.repo.On("FindByUUID", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Вызов
	result, err := s.service.PayOrder(ctx, orderUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrOrderNotFound, err)
	s.Empty(result)

	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_PaymentError() {
	// Подготовка
	ctx := context.Background()
	orderUUID := uuid.New().String()
	paymentMethod := "CARD"
	expectedErr := errors.New("payment service error")

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New().String(),
		PartUUIDs:  []string{uuid.New().String()},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	s.repo.On("FindByUUID", ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", ctx, order.UserUUID, orderUUID, paymentMethod).Return("", expectedErr)

	// Вызов
	result, err := s.service.PayOrder(ctx, orderUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Contains(err.Error(), "payment failed")
	s.Empty(result)

	s.repo.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}
