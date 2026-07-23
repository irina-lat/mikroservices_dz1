package v1

import (
	orderapi "shared/pkg/openapi/order/v1"

	"order/internal/service/order"
)

// API реализует HTTP хендлеры для OrderService
type API struct {
	orderService order.Service
}

// Убеждаемся, что API реализует интерфейс orderapi.Handler
var _ orderapi.Handler = (*API)(nil)
