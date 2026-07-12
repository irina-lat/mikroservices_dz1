//go:generate mockery --name=Service --output=../mocks --case=underscore
//go:generate mockery --name=Repository --output=../mocks --case=underscore
//go:generate mockery --name=InventoryClient --output=../mocks --case=underscore
//go:generate mockery --name=PaymentClient --output=../mocks --case=underscore

package service

import (
	"github.com/irina-lat/microservices-course/order/internal/service/order"
)

// OrderService - алиас для интерфейса сервиса
type OrderService = order.Service
