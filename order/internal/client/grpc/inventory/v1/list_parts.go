package v1

import (
	"context"

	"order/internal/client/converter"
	"order/internal/model"
	"platform/pkg/middleware/grpc"
	inventorypb "shared/pkg/proto/inventory/v1"
)

func (c *InventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]*model.Part, error) {
	// Передаём session-uuid в gRPC metadata
	ctx = grpc.ForwardSessionUUIDToGRPC(ctx)

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