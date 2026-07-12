package v1

import (
	"github.com/irina-lat/microservices-course/inventory/internal/service/part"
	pb "shared/pkg/proto/inventory/v1"
)

// API реализует gRPC хендлеры для InventoryService
type API struct {
	pb.UnimplementedInventoryServiceServer
	service part.Service
}

// NewAPI создаёт новый экземпляр API
func NewAPI(service part.Service) *API {
	return &API{
		service: service,
	}
}
