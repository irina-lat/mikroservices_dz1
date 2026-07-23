package v1

import (
	"payment/internal/service/payment"
	pb "shared/pkg/proto/payment/v1"
)

// API реализует gRPC хендлеры для PaymentService
type API struct {
	pb.UnimplementedPaymentServiceServer
	service payment.Service
}

// NewAPI создаёт новый экземпляр API
func NewAPI(service payment.Service) *API {
	return &API{
		service: service,
	}
}
