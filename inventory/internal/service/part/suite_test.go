package part

import (
	"testing"

	"github.com/irina-lat/microservices-course/inventory/internal/repository/mocks"


	"github.com/stretchr/testify/suite"
)

type PartServiceTestSuite struct {
	suite.Suite
	service *PartService
	repo    *mocks.MockRepository
}


func (s *PartServiceTestSuite) SetupTest() {
	// Создаём мок репозитория
	s.repo = mocks.NewMockRepository(s.T())

	// Создаём сервис с моком
	s.service = NewService(s.repo)
}

func TestPartService(t *testing.T) {
	suite.Run(t, new(PartServiceTestSuite))
}
