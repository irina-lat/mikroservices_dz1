package handler

import (
	"context"

	"inventory/internal/service"
	"inventory/pkg/model"
	pb "shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedInventoryServiceServer
	partService *service.PartService
}

func NewGrpcHandler(partService *service.PartService) *GrpcHandler {
	return &GrpcHandler{
		partService: partService,
	}
}

// GetPart реализует gRPC метод GetPart
func (h *GrpcHandler) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	part, err := h.partService.GetPart(ctx, req.Uuid)
	if err != nil {
		if err.Error() == "part not found" {
			return nil, status.Error(codes.NotFound, "part not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetPartResponse{
		Part: h.convertToProtoPart(part),
	}, nil
}

// ListParts реализует gRPC метод ListParts
func (h *GrpcHandler) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	// Конвертируем proto фильтр в модель
	filter := h.convertToModelFilter(req.Filter)

	parts, err := h.partService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoParts := make([]*pb.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, h.convertToProtoPart(part))
	}

	return &pb.ListPartsResponse{
		Parts: protoParts,
	}, nil
}

// convertToProtoPart конвертирует модель Part в proto Part
func (h *GrpcHandler) convertToProtoPart(part *model.Part) *pb.Part {
	return &pb.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      pb.Category(part.Category),
		Dimensions: &pb.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &pb.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  h.convertToProtoMetadata(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// convertToModelFilter конвертирует proto фильтр в модель фильтра
func (h *GrpcHandler) convertToModelFilter(protoFilter *pb.PartsFilter) *model.PartsFilter {
	if protoFilter == nil {
		return nil
	}

	filter := &model.PartsFilter{
		UUIDs:  protoFilter.Uuids,
		Names:  protoFilter.Names,
		Tags:   protoFilter.Tags,
		ManufacturerCountries: protoFilter.ManufacturerCountries,
	}

	// Конвертируем категории
	for _, cat := range protoFilter.Categories {
		filter.Categories = append(filter.Categories, model.Category(cat))
	}

	return filter
}

// convertToProtoMetadata конвертирует метаданные
func (h *GrpcHandler) convertToProtoMetadata(metadata map[string]*model.Value) map[string]*pb.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*pb.Value)
	for key, value := range metadata {
		pbValue := &pb.Value{}
		if value.StringValue != nil {
			pbValue.Kind = &pb.Value_StringValue{StringValue: *value.StringValue}
		} else if value.Int64Value != nil {
			pbValue.Kind = &pb.Value_Int64Value{Int64Value: *value.Int64Value}
		} else if value.DoubleValue != nil {
			pbValue.Kind = &pb.Value_DoubleValue{DoubleValue: *value.DoubleValue}
		} else if value.BoolValue != nil {
			pbValue.Kind = &pb.Value_BoolValue{BoolValue: *value.BoolValue}
		}
		result[key] = pbValue
	}
	return result
}