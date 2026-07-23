package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"order/internal/model"
	orderapi "shared/pkg/openapi/order/v1"
)

// CreateOrder обрабатывает POST /api/v1/orders
func (a *API) CreateOrder(ctx context.Context, req *orderapi.CreateOrderRequest) (orderapi.CreateOrderRes, error) {
	log.Printf("CreateOrder: user_uuid=%s, part_uuids=%v", req.UserUUID, req.PartUuids)

	partUUIDs := make([]string, len(req.PartUuids))
	for i, u := range req.PartUuids {
		partUUIDs[i] = u.String()
	}

	orderModel, err := a.orderService.CreateOrder(ctx, req.UserUUID.String(), partUUIDs)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderapi.BadRequest{
				Message: orderapi.NewOptString("Some parts not found"),
				Code:    orderapi.NewOptString("PART_NOT_FOUND"),
			}, nil
		}
		log.Printf("CreateOrder error: %v", err)
		return &orderapi.InternalServerError{
			Message: orderapi.NewOptString("Internal server error"),
			Code:    orderapi.NewOptString("INTERNAL_ERROR"),
		}, nil
	}

	return &orderapi.CreateOrderResponse{
		OrderUUID:  uuid.MustParse(orderModel.OrderUUID),
		TotalPrice: orderModel.TotalPrice,
	}, nil
}
