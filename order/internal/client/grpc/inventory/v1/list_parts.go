package v1

import (
	"context"

	"github.com/irina-lat/microservices-course/order/internal/client/converter"
	"github.com/irina-lat/microservices-course/order/internal/model"
	inventorypb "shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей по UUID
func (c *InventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]*model.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: partUUIDs,
		},
	})

	if err != nil {
		return nil, err
	}

	return converter.ProtoPartsToModels(resp.Parts), nil
}
