package payment

import (
	"context"

	"payment/internal/model"
)

func (s *PaymentServiceTestSuite) TestPayOrder_Success() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := "test-user-uuid-456"
	paymentMethod := "CARD"

	// Вызов
	transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.NoError(err)
	s.NotEmpty(transactionUUID)
	s.NotEqual(orderUUID, transactionUUID) // transactionUUID должен быть новым
	s.NotEqual(userUUID, transactionUUID)
}

func (s *PaymentServiceTestSuite) TestPayOrder_EmptyOrderUUID() {
	// Подготовка
	ctx := context.Background()
	orderUUID := ""
	userUUID := "test-user-uuid-456"
	paymentMethod := "CARD"

	// Вызов
	transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrEmptyOrderUUID, err)
	s.Empty(transactionUUID)
}

func (s *PaymentServiceTestSuite) TestPayOrder_EmptyUserUUID() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := ""
	paymentMethod := "CARD"

	// Вызов
	transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrEmptyUserUUID, err)
	s.Empty(transactionUUID)
}

func (s *PaymentServiceTestSuite) TestPayOrder_EmptyPaymentMethod() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := "test-user-uuid-456"
	paymentMethod := ""

	// Вызов
	transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrInvalidPaymentMethod, err)
	s.Empty(transactionUUID)
}

func (s *PaymentServiceTestSuite) TestPayOrder_InvalidPaymentMethod() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := "test-user-uuid-456"
	paymentMethod := "INVALID_METHOD"

	// Вызов
	transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.Error(err)
	s.Equal(model.ErrInvalidPaymentMethod, err)
	s.Empty(transactionUUID)
}

func (s *PaymentServiceTestSuite) TestPayOrder_AllPaymentMethods() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := "test-user-uuid-456"
	paymentMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}

	for _, method := range paymentMethods {
		// Вызов
		transactionUUID, err := s.service.PayOrder(ctx, orderUUID, userUUID, method)

		// Проверки
		s.NoError(err)
		s.NotEmpty(transactionUUID)
		s.NotEqual(orderUUID, transactionUUID)
	}
}

func (s *PaymentServiceTestSuite) TestPayOrder_GeneratesUniqueUUID() {
	// Подготовка
	ctx := context.Background()
	orderUUID := "test-order-uuid-123"
	userUUID := "test-user-uuid-456"
	paymentMethod := "CARD"

	// Вызов несколько раз
	uuid1, err1 := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)
	uuid2, err2 := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Проверки
	s.NoError(err1)
	s.NoError(err2)
	s.NotEmpty(uuid1)
	s.NotEmpty(uuid2)
	s.NotEqual(uuid1, uuid2) // UUID должны быть разными
}
