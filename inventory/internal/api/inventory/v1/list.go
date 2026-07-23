package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory/internal/converter"
	pb "shared/pkg/proto/inventory/v1"
)

// ListParts обрабатывает gRPC запрос ListParts
func (a *API) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	filter := converter.ProtoFilterToModel(req.Filter)

	parts, err := a.service.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListPartsResponse{
		Parts: converter.ModelPartsToProto(parts),
	}, nil
}
