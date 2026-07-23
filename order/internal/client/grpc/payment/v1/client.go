//go:generate mockery --name=PaymentClient --output=../../../mocks --case=underscore

package v1

import (
	paymentpb "shared/pkg/proto/payment/v1"
)

// PaymentClient представляет клиент для PaymentService
type PaymentClient struct {
	client paymentpb.PaymentServiceClient
}

// NewPaymentClient создаёт новый клиент для PaymentService
func NewPaymentClient(client paymentpb.PaymentServiceClient) *PaymentClient {
	return &PaymentClient{
		client: client,
	}
}
