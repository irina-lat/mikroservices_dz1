package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"order/internal/model"
	"platform/pkg/middleware/grpc"
)

func (s *OrderService) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	// Передаём session-uuid в gRPC вызов InventoryService
	ctx = grpc.ForwardSessionUUIDToGRPC(ctx)

	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}

	if len(parts) != len(partUUIDs) {
		return nil, model.ErrPartNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	orderUUID := uuid.New().String()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}