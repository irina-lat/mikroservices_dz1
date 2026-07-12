//go:generate mockery --name=Repository --output=../mocks --case=underscore

package repository

import (
	"github.com/irina-lat/microservices-course/order/internal/repository/order"
)

// OrderRepository - алиас для интерфейса репозитория
type OrderRepository = order.Repository
