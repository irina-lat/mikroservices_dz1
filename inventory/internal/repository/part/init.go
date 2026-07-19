package part

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/irina-lat/microservices-course/inventory/internal/model"
)

// InitSampleData заполняет MongoDB тестовыми данными
func (r *MongoRepository) InitSampleData(ctx context.Context) error {
	now := time.Now()

	sampleParts := []*model.Part{
		{
			UUID:          uuid.New().String(),
			Name:          "Main Engine",
			Description:   "Main propulsion engine for spacecraft",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions: model.Dimensions{
				Length: 200.0,
				Width:  150.0,
				Height: 100.0,
				Weight: 500.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "https://spacex.com",
			},
			Tags:      []string{"engine", "propulsion", "main"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Fuel Tank",
			Description:   "High-capacity fuel tank",
			Price:         5000.0,
			StockQuantity: 10,
			Category:      model.CategoryFuel,
			Dimensions: model.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 200.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Boeing",
				Country: "USA",
				Website: "https://boeing.com",
			},
			Tags:      []string{"fuel", "tank", "storage"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Porthole Window",
			Description:   "Reinforced viewing window",
			Price:         2500.0,
			StockQuantity: 20,
			Category:      model.CategoryPorthole,
			Dimensions: model.Dimensions{
				Length: 50.0,
				Width:  50.0,
				Height: 10.0,
				Weight: 25.0,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Thales",
				Country: "France",
				Website: "https://thalesgroup.com",
			},
			Tags:      []string{"window", "viewing", "porthole"},
			Metadata:  make(map[string]*model.Value),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, part := range sampleParts {
		if err := r.Save(ctx, part); err != nil {
			return err
		}
	}

	return nil
}