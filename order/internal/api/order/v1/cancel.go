package v1

import (
	"context"
	"errors"
	"log"

	"github.com/irina-lat/microservices-course/order/internal/model"
	orderapi "shared/pkg/openapi/order/v1"
)

// CancelOrder обрабатывает POST /api/v1/orders/{order_uuid}/cancel
func (a *API) CancelOrder(ctx context.Context, params orderapi.CancelOrderParams) (orderapi.CancelOrderRes, error) {
	log.Printf("CancelOrder: order_uuid=%s", params.OrderUUID)

	err := a.orderService.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderapi.NotFound{
				Message: orderapi.NewOptString("Order not found"),
				Code:    orderapi.NewOptString("NOT_FOUND"),
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderapi.Conflict{
				Message: orderapi.NewOptString("Cannot cancel paid order"),
				Code:    orderapi.NewOptString("ORDER_ALREADY_PAID"),
			}, nil
		}
		log.Printf("CancelOrder error: %v", err)
		return &orderapi.InternalServerError{
			Message: orderapi.NewOptString("Internal server error"),
			Code:    orderapi.NewOptString("INTERNAL_ERROR"),
		}, nil
	}

	return &orderapi.CancelOrderNoContent{}, nil
}

