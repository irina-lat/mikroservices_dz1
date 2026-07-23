package v1

import (
	"order/internal/service/order"
)

// NewAPI создаёт новый экземпляр API
func NewAPI(orderService order.Service) *API {
	return &API{
		orderService: orderService,
	}
}
