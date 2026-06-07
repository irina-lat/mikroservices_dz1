package handler

import (
	"context"

	"payment/internal/service"
	"payment/pkg/model"
	pb "shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *service.PaymentService
}

func NewGrpcHandler(paymentService *service.PaymentService) *GrpcHandler {
	return &GrpcHandler{
		paymentService: paymentService,
	}
}

// PayOrder реализует gRPC метод PayOrder
func (h *GrpcHandler) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	// Валидация обязательных полей
	if req.OrderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}
	if req.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}
	if req.PaymentMethod == pb.PaymentMethod_PAYMENT_METHOD_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "payment_method is required or unknown")
	}

	// Конвертируем proto запрос в модель
	serviceReq := &model.PayOrderRequest{
		OrderUUID:     req.OrderUuid,
		UserUUID:      req.UserUuid,
		PaymentMethod: h.convertToModelPaymentMethod(req.PaymentMethod),
	}

	// Вызываем сервис
	resp, err := h.paymentService.PayOrder(ctx, serviceReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Возвращаем proto ответ
	return &pb.PayOrderResponse{
		TransactionUuid: resp.TransactionUUID,
	}, nil
}

// convertToModelPaymentMethod конвертирует proto PaymentMethod в модель
func (h *GrpcHandler) convertToModelPaymentMethod(method pb.PaymentMethod) model.PaymentMethod {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case pb.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case pb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}