package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"payment/internal/model"
)

// PayOrder обрабатывает оплату заказа
func (s *PaymentService) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	// Валидация
	if orderUUID == "" {
		return "", model.ErrEmptyOrderUUID
	}
	if userUUID == "" {
		return "", model.ErrEmptyUserUUID
	}

	// Проверка метода оплаты
	isValidMethod := false
	validMethods := []string{"CARD", "SBP", "CREDIT_CARD", "INVESTOR_MONEY"}
	for _, m := range validMethods {
		if m == paymentMethod {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		return "", model.ErrInvalidPaymentMethod
	}


	// Генерируем transaction_uuid
	transactionUUID := uuid.New().String()

	// Выводим сообщение в консоль
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)
	log.Printf("Детали оплаты: order_uuid=%s, user_uuid=%s, payment_method=%s",
		orderUUID, userUUID, paymentMethod)

	return transactionUUID, nil
}
