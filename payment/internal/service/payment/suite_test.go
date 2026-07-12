package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	service *PaymentService
}


func (s *PaymentServiceTestSuite) SetupTest() {
	// Создаём сервис (без моков, так как он не зависит от внешних зависимостей)
	s.service = New()
}

func TestPaymentService(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
