package v1

import (
	"github.com/irina-lat/microservices-course/order/internal/service/order"
	orderapi "shared/pkg/openapi/order/v1"
)

// API реализует HTTP хендлеры для OrderService
type API struct {
	orderService order.Service
}

// Убеждаемся, что API реализует интерфейс orderapi.Handler
var _ orderapi.Handler = (*API)(nil)
