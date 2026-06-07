package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"payment/pkg/model"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// PayOrder обрабатывает оплату заказа
func (s *PaymentService) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	// Генерируем transaction_uuid
	transactionUUID := uuid.New().String()

	// Выводим сообщение в консоль
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)
	log.Printf("Детали оплаты: order_uuid=%s, user_uuid=%s, payment_method=%s",
		req.OrderUUID, req.UserUUID, req.PaymentMethod)

	return &model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}