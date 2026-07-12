package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/irina-lat/microservices-course/order/internal/model"
)

func (s *OrderServiceTestSuite) TestCreateOrder_Success() {
	// Подготовка
	ctx := context.Background()
	userUUID := uuid.New().String()
	partUUIDs := []string{uuid.New().String(), uuid.New().String()}

	expectedParts := []*model.Part{
		{UUID: partUUIDs[0], Price: 100.0},
		{UUID: partUUIDs[1], Price: 200.0},
	}

	// Настройка моков
	s.inventoryClient.On("ListParts", ctx, partUUIDs).Return(expectedParts, nil)
	s.repo.On("Save", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	// Вызов
	order, err := s.service.CreateOrder(ctx, userUUID, partUUIDs)

	// Проверки
	s.NoError(err)
	s.NotNil(order)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(partUUIDs, order.PartUUIDs)
	s.Equal(300.0, order.TotalPrice)
	s.Equal(model.StatusPendingPayment, order.Status)

	s.inventoryClient.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCreateOrder_PartNotFound() {
	// Подготовка
	ctx := context.Background()
	userUUID := uuid.New().String()
	partUUIDs := []string{uuid.New().String()}

	// Настройка моков - возвращаем пустой список (деталь не найдена)
	s.inventoryClient.On("ListParts", ctx, partUUIDs).Return([]*model.Part{}, nil)

	// Вызов
	order, err := s.service.CreateOrder(ctx, userUUID, partUUIDs)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrPartNotFound, err)
	s.Nil(order)

	s.inventoryClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCreateOrder_InventoryError() {
	// Подготовка
	ctx := context.Background()
	userUUID := uuid.New().String()
	partUUIDs := []string{uuid.New().String()}
	expectedErr := errors.New("inventory service error")

	// Настройка моков - InventoryService возвращает ошибку
	s.inventoryClient.On("ListParts", ctx, partUUIDs).Return(nil, expectedErr)

	// Вызов
	order, err := s.service.CreateOrder(ctx, userUUID, partUUIDs)

	// Проверки
	s.Error(err)
	s.Contains(err.Error(), "failed to get parts")
	s.Nil(order)

	s.inventoryClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCreateOrder_RepositoryError() {
	// Подготовка
	ctx := context.Background()
	userUUID := uuid.New().String()
	partUUIDs := []string{uuid.New().String()}

	expectedParts := []*model.Part{
		{UUID: partUUIDs[0], Price: 100.0},
	}
	expectedErr := errors.New("repository error")

	// Настройка моков
	s.inventoryClient.On("ListParts", ctx, partUUIDs).Return(expectedParts, nil)
	s.repo.On("Save", ctx, mock.AnythingOfType("*model.Order")).Return(expectedErr)

	// Вызов
	order, err := s.service.CreateOrder(ctx, userUUID, partUUIDs)

	// Проверки
	s.Error(err)
	s.Contains(err.Error(), "failed to save order")
	s.Nil(order)

	s.inventoryClient.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}
