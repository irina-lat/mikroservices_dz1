package converter

import (
	"inventory/internal/model"
	repomodel "inventory/internal/repository/model"
)

func ToRepoPart(part *model.Part) *repomodel.Part {
	if part == nil {
		return nil
	}

	return &repomodel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repomodel.Category(part.Category),
		Dimensions: repomodel.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: repomodel.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  convertMetadataToRepo(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ToServicePart(part *repomodel.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions: model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  convertMetadataToService(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ToServiceParts(parts []*repomodel.Part) []*model.Part {
	if parts == nil {
		return nil
	}

	result := make([]*model.Part, len(parts))
	for i, part := range parts {
		result[i] = ToServicePart(part)
	}
	return result
}

func convertMetadataToRepo(metadata map[string]*model.Value) map[string]*repomodel.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*repomodel.Value)
	for key, value := range metadata {
		repoValue := &repomodel.Value{}
		if value.StringValue != nil {
			repoValue.StringValue = value.StringValue
		}
		if value.Int64Value != nil {
			repoValue.Int64Value = value.Int64Value
		}
		if value.DoubleValue != nil {
			repoValue.DoubleValue = value.DoubleValue
		}
		if value.BoolValue != nil {
			repoValue.BoolValue = value.BoolValue
		}
		result[key] = repoValue
	}
	return result
}

func convertMetadataToService(metadata map[string]*repomodel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*model.Value)
	for key, value := range metadata {
		serviceValue := &model.Value{}
		if value.StringValue != nil {
			serviceValue.StringValue = value.StringValue
		}
		if value.Int64Value != nil {
			serviceValue.Int64Value = value.Int64Value
		}
		if value.DoubleValue != nil {
			serviceValue.DoubleValue = value.DoubleValue
		}
		if value.BoolValue != nil {
			serviceValue.BoolValue = value.BoolValue
		}
		result[key] = serviceValue
	}
	return result
}
