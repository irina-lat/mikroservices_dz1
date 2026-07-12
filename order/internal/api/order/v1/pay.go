package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/irina-lat/microservices-course/order/internal/model"
	orderapi "shared/pkg/openapi/order/v1"
)

// PayOrder обрабатывает POST /api/v1/orders/{order_uuid}/pay
func (a *API) PayOrder(ctx context.Context, req *orderapi.PayOrderRequest, params orderapi.PayOrderParams) (orderapi.PayOrderRes, error) {
	log.Printf("PayOrder: order_uuid=%s, payment_method=%s", params.OrderUUID, req.PaymentMethod)

	transactionUUID, err := a.orderService.PayOrder(ctx, params.OrderUUID.String(), string(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderapi.NotFound{
				Message: orderapi.NewOptString("Order not found"),
				Code:    orderapi.NewOptString("NOT_FOUND"),
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderapi.Conflict{
				Message: orderapi.NewOptString("Order already paid"),
				Code:    orderapi.NewOptString("ALREADY_PAID"),
			}, nil
		}
		log.Printf("PayOrder error: %v", err)
		return &orderapi.InternalServerError{
			Message: orderapi.NewOptString("Internal server error"),
			Code:    orderapi.NewOptString("INTERNAL_ERROR"),
		}, nil
	}

	return &orderapi.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID),
	}, nil
}

