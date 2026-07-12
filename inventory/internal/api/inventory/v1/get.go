package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/irina-lat/microservices-course/inventory/internal/converter"
	"github.com/irina-lat/microservices-course/inventory/internal/model"
	pb "shared/pkg/proto/inventory/v1"
)

// GetPart обрабатывает gRPC запрос GetPart
func (a *API) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	part, err := a.service.GetPart(ctx, req.Uuid)
	if err != nil {
		if err == model.ErrPartNotFound {
			return nil, status.Error(codes.NotFound, "part not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetPartResponse{
		Part: converter.ModelPartToProto(part),
	}, nil
}
