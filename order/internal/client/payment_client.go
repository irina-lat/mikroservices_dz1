package client

import (
	"context"

	pb "shared/pkg/proto/payment/v1"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (string, error)
}

type GrpcPaymentClient struct {
	client pb.PaymentServiceClient
}

func NewGrpcPaymentClient(client pb.PaymentServiceClient) *GrpcPaymentClient {
	return &GrpcPaymentClient{client: client}
}

func (c *GrpcPaymentClient) PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (string, error) {
	// Конвертируем строку в enum
	var method pb.PaymentMethod
	switch paymentMethod {
	case "CARD":
		method = pb.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		method = pb.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		method = pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY":
		method = pb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		method = pb.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}

	resp, err := c.client.PayOrder(ctx, &pb.PayOrderRequest{
		UserUuid:      userUUID,
		OrderUuid:     orderUUID,
		PaymentMethod: method,
	})
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}