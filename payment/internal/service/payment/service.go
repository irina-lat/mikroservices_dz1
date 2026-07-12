package payment

import (
	"context"
)

// Service определяет интерфейс бизнес-логики для оплаты
type Service interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

// PaymentService реализует бизнес-логику
type PaymentService struct{}

// New создаёт новый экземпляр PaymentService
func New() *PaymentService {
	return &PaymentService{}
}