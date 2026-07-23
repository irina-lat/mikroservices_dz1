package order

import (
	"context"
	"fmt"

	"order/internal/model"
)

// PayOrder оплачивает заказ
func (s *OrderService) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (string, error) {
	// 1. Находим заказ
	order, err := s.repo.FindByUUID(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	// 2. Проверяем статус
	if order.Status == model.StatusPaid {
		return "", model.ErrOrderAlreadyPaid
	}
	if order.Status == model.StatusCancelled {
		return "", model.ErrOrderAlreadyCanceled
	}

	// 3. Вызываем PaymentService
	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		return "", fmt.Errorf("payment failed: %w", err)
	}

	// 4. Обновляем заказ
	method := model.PaymentMethod(paymentMethod)
	order.Status = model.StatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method

	if err := s.repo.Update(ctx, order); err != nil {
		return "", fmt.Errorf("failed to update order: %w", err)
	}

	return transactionUUID, nil
}
