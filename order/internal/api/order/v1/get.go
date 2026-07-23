package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"order/internal/model"
	orderapi "shared/pkg/openapi/order/v1"
)

// GetOrder обрабатывает GET /api/v1/orders/{order_uuid}
func (a *API) GetOrder(ctx context.Context, params orderapi.GetOrderParams) (orderapi.GetOrderRes, error) {
	log.Printf("GetOrder: order_uuid=%s", params.OrderUUID)

	orderModel, err := a.orderService.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderapi.NotFound{
				Message: orderapi.NewOptString("Order not found"),
				Code:    orderapi.NewOptString("NOT_FOUND"),
			}, nil
		}
		log.Printf("GetOrder error: %v", err)
		return &orderapi.InternalServerError{
			Message: orderapi.NewOptString("Internal server error"),
			Code:    orderapi.NewOptString("INTERNAL_ERROR"),
		}, nil
	}

	return a.convertToOrderResponse(orderModel), nil
}

// convertToOrderResponse преобразует модель в API ответ
func (a *API) convertToOrderResponse(order *model.Order) *orderapi.Order {
	resp := &orderapi.Order{
		OrderUUID:  uuid.MustParse(order.OrderUUID),
		UserUUID:   uuid.MustParse(order.UserUUID),
		PartUuids:  make([]uuid.UUID, len(order.PartUUIDs)),
		TotalPrice: order.TotalPrice,
		Status:     orderapi.OrderStatus(order.Status),
	}

	for i, partUUID := range order.PartUUIDs {
		resp.PartUuids[i] = uuid.MustParse(partUUID)
	}

	if order.TransactionUUID != nil {
		transactionUUID := uuid.MustParse(*order.TransactionUUID)
		resp.TransactionUUID = orderapi.NewOptNilUUID(transactionUUID)
	}

	if order.PaymentMethod != nil {
		paymentMethod := orderapi.OrderPaymentMethod(*order.PaymentMethod)
		resp.PaymentMethod = orderapi.NewOptNilOrderPaymentMethod(paymentMethod)
	}

	return resp
}
