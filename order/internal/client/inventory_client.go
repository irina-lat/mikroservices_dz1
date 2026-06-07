package client

import (
	"context"

	pb "shared/pkg/proto/inventory/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]*pb.Part, error)
}

type GrpcInventoryClient struct {
	client pb.InventoryServiceClient
}

func NewGrpcInventoryClient(client pb.InventoryServiceClient) *GrpcInventoryClient {
	return &GrpcInventoryClient{client: client}
}

func (c *GrpcInventoryClient) ListParts(ctx context.Context, partUUIDs []string) ([]*pb.Part, error) {
	resp, err := c.client.ListParts(ctx, &pb.ListPartsRequest{
		PartUuids: partUUIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}