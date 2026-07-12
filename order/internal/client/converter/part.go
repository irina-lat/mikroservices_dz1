package converter

import (
	"github.com/irina-lat/microservices-course/order/internal/model"
	inventorypb "shared/pkg/proto/inventory/v1"
)

// ProtoPartToModel конвертирует proto Part в модель Part
func ProtoPartToModel(part *inventorypb.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category.String(),
	}
}

// ProtoPartsToModels конвертирует список proto Part в список моделей
func ProtoPartsToModels(parts []*inventorypb.Part) []*model.Part {
	if parts == nil {
		return nil
	}

	result := make([]*model.Part, len(parts))
	for i, part := range parts {
		result[i] = ProtoPartToModel(part)
	}
	return result
}
