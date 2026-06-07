package handler

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"order/internal/repository"
	"order/internal/service"
	orderapi "shared/pkg/openapi/order/v1"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder обрабатывает POST /api/v1/orders
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderapi.CreateOrderReq) (*orderapi.CreateOrderResponse, error) {
	log.Printf("CreateOrder: user_uuid=%s, part_uuids=%v", req.UserUUID, req.PartUUIDs)

	order, err := h.orderService.CreateOrder(ctx, req.UserUUID.String(), req.PartUUIDs)
	if err != nil {
		if errors.Is(err, service.ErrPartNotFound) {
			return nil, &orderapi.BadRequestError{
				Message: "Some parts not found",
				Code:    "PART_NOT_FOUND",
			}
		}
		log.Printf("CreateOrder error: %v", err)
		return nil, &orderapi.InternalServerError{
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
		}
	}

	return &orderapi.CreateOrderResponse{
		OrderUUID:  uuid.MustParse(order.OrderUUID),
		TotalPrice: order.TotalPrice,
	}, nil
}

// GetOrder обрабатывает GET /api/v1/orders/{order_uuid}
func (h *OrderHandler) GetOrder(ctx context.Context, params orderapi.GetOrderParams) (*orderapi.Order, error) {
	log.Printf("GetOrder: order_uuid=%s", params.OrderUUID)

	order, err := h.orderService.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, &orderapi.NotFoundError{
				Message: "Order not found",
				Code:    "NOT_FOUND",
			}
		}
		log.Printf("GetOrder error: %v", err)
		return nil, &orderapi.InternalServerError{
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
		}
	}

	return h.mapToOrderResponse(order), nil
}

// PayOrder обрабатывает POST /api/v1/orders/{order_uuid}/pay
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderapi.PayOrderReq, params orderapi.PayOrderParams) (*orderapi.PayOrderResponse, error) {
	log.Printf("PayOrder: order_uuid=%s, payment_method=%s", params.OrderUUID, req.PaymentMethod)

	transactionUUID, err := h.orderService.PayOrder(ctx, params.OrderUUID.String(), string(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, &orderapi.NotFoundError{
				Message: "Order not found",
				Code:    "NOT_FOUND",
			}
		}
		if errors.Is(err, service.ErrOrderAlreadyPaid) {
			return nil, &orderapi.ConflictError{
				Message: "Order already paid",
				Code:    "ALREADY_PAID",
			}
		}
		log.Printf("PayOrder error: %v", err)
		return nil, &orderapi.InternalServerError{
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
		}
	}

	return &orderapi.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}

// CancelOrder обрабатывает POST /api/v1/orders/{order_uuid}/cancel
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderapi.CancelOrderParams) error {
	log.Printf("CancelOrder: order_uuid=%s", params.OrderUUID)

	err := h.orderService.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return &orderapi.NotFoundError{
				Message: "Order not found",
				Code:    "NOT_FOUND",
			}
		}
		if errors.Is(err, service.ErrOrderAlreadyPaid) {
			return &orderapi.ConflictError{
				Message: "Cannot cancel paid order",
				Code:    "ORDER_ALREADY_PAID",
			}
		}
		log.Printf("CancelOrder error: %v", err)
		return &orderapi.InternalServerError{
			Message: "Internal server error",
			Code:    "INTERNAL_ERROR",
		}
	}

	return nil
}

// mapToOrderResponse преобразует модель в API ответ
func (h *OrderHandler) mapToOrderResponse(order *service.Order) *orderapi.Order {
	resp := &orderapi.Order{
		OrderUUID:  uuid.MustParse(order.OrderUUID),
		UserUUID:   uuid.MustParse(order.UserUUID),
		PartUUIDs:  make([]uuid.UUID, len(order.PartUUIDs)),
		TotalPrice: order.TotalPrice,
		Status:     orderapi.OrderStatus(order.Status),
	}

	for i, partUUID := range order.PartUUIDs {
		resp.PartUUIDs[i] = uuid.MustParse(partUUID)
	}

	if order.TransactionUUID != nil {
		transactionUUID := uuid.MustParse(*order.TransactionUUID)
		resp.TransactionUUID = &transactionUUID
	}

	if order.PaymentMethod != nil {
		paymentMethod := orderapi.PayOrderReqPaymentMethod(*order.PaymentMethod)
		resp.PaymentMethod = &paymentMethod
	}

	return resp
}