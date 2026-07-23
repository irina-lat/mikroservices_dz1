package v1

import (
	"context"

	paymentpb "shared/pkg/proto/payment/v1"
)

// PayOrder оплачивает заказ
func (c *PaymentClient) PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (string, error) {
	// Конвертируем строку в proto enum
	var method paymentpb.PaymentMethod
	switch paymentMethod {
	case "CARD":
		method = paymentpb.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		method = paymentpb.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		method = paymentpb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY":
		method = paymentpb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		method = paymentpb.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}

	resp, err := c.client.PayOrder(ctx, &paymentpb.PayOrderRequest{
		UserUuid:      userUUID,
		OrderUuid:     orderUUID,
		PaymentMethod: method,
	})
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
