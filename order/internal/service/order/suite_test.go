package order

import (
	"testing"

	"github.com/irina-lat/microservices-course/order/internal/service/mocks"

	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	suite.Suite
	service         *OrderService
	repo            *mocks.MockRepository
	inventoryClient *mocks.MockInventoryClient
	paymentClient   *mocks.MockPaymentClient
}

func (s *OrderServiceTestSuite) SetupTest() {
	// Создаём моки
	s.repo = mocks.NewMockRepository(s.T())
	s.inventoryClient = mocks.NewMockInventoryClient(s.T())
	s.paymentClient = mocks.NewMockPaymentClient(s.T())

	// Создаём сервис с моками
	s.service = NewService(s.repo, s.inventoryClient, s.paymentClient)
}

func TestOrderService(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
