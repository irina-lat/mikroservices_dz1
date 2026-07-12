package v1

import (
	"github.com/irina-lat/microservices-course/order/internal/service/order"
)

// NewAPI создаёт новый экземпляр API
func NewAPI(orderService order.Service) *API {
	return &API{
		orderService: orderService,
	}
}
