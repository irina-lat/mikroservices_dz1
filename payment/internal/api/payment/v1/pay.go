package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment/internal/model"
	pb "shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает gRPC запрос PayOrder
func (a *API) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	// Валидация
	if req.OrderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}
	if req.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	// Конвертируем proto в модель
	paymentMethod := convertProtoToModelPaymentMethod(req.PaymentMethod)
	if paymentMethod == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_method is required")
	}

	// Вызываем сервис
	transactionUUID, err := a.service.PayOrder(ctx, req.OrderUuid, req.UserUuid, paymentMethod)
	if err != nil {
		switch err {
		case model.ErrInvalidPaymentMethod:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case model.ErrEmptyOrderUUID:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case model.ErrEmptyUserUUID:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

// convertProtoToModelPaymentMethod конвертирует proto PaymentMethod в строку
func convertProtoToModelPaymentMethod(method pb.PaymentMethod) string {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_CARD:
		return "CARD"
	case pb.PaymentMethod_PAYMENT_METHOD_SBP:
		return "SBP"
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return "CREDIT_CARD"
	case pb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return ""
	}
}
