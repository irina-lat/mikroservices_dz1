//go:generate mockery --name=InventoryClient --output=../../../mocks --case=underscore

package v1

import (
	inventorypb "shared/pkg/proto/inventory/v1"
)

// InventoryClient представляет клиент для InventoryService
type InventoryClient struct {
	client inventorypb.InventoryServiceClient
}

// NewInventoryClient создаёт новый клиент для InventoryService
func NewInventoryClient(client inventorypb.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{
		client: client,
	}
}