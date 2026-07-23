package converter

import (
	"time"

	"inventory/internal/model"
	pb "shared/pkg/proto/inventory/v1"
)

// ModelPartToProto конвертирует модель Part в proto Part
func ModelPartToProto(part *model.Part) *pb.Part {
	if part == nil {
		return nil
	}

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
		Metadata:  ModelMetadataToProto(part.Metadata),
		CreatedAt: part.CreatedAt.Format(time.RFC3339),
		UpdatedAt: part.UpdatedAt.Format(time.RFC3339),
	}
}

// ModelPartsToProto конвертирует срез моделей в срез proto
func ModelPartsToProto(parts []*model.Part) []*pb.Part {
	if parts == nil {
		return nil
	}

	result := make([]*pb.Part, len(parts))
	for i, part := range parts {
		result[i] = ModelPartToProto(part)
	}
	return result
}

// ProtoFilterToModel конвертирует proto фильтр в модель фильтра
func ProtoFilterToModel(filter *pb.PartsFilter) *model.PartsFilter {
	if filter == nil {
		return nil
	}

	return &model.PartsFilter{
		UUIDs:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            ProtoCategoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// ProtoCategoriesToModel конвертирует срез proto категорий в модель
func ProtoCategoriesToModel(categories []pb.Category) []model.Category {
	if categories == nil {
		return nil
	}

	result := make([]model.Category, len(categories))
	for i, cat := range categories {
		result[i] = model.Category(cat)
	}
	return result
}

// ModelMetadataToProto конвертирует модель метаданных в proto
func ModelMetadataToProto(metadata map[string]*model.Value) map[string]*pb.Value {
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
