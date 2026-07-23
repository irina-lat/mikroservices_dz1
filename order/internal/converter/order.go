package converter

import (
	"order/internal/model"
	repomodel "order/internal/repository/model"
)

// ToRepoOrder конвертирует модель сервисного слоя в модель репозитория
func ToRepoOrder(order *model.Order) *repomodel.Order {
	if order == nil {
		return nil
	}

	return &repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*repomodel.PaymentMethod)(order.PaymentMethod),
		Status:          repomodel.OrderStatus(order.Status),
	}
}

// ToServiceOrder конвертирует модель репозитория в модель сервисного слоя
func ToServiceOrder(order *repomodel.Order) *model.Order {
	if order == nil {
		return nil
	}

	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
	}
}

// ToServiceOrders конвертирует срез моделей репозитория в срез моделей сервисного слоя
func ToServiceOrders(orders []*repomodel.Order) []*model.Order {
	if orders == nil {
		return nil
	}

	result := make([]*model.Order, len(orders))
	for i, order := range orders {
		result[i] = ToServiceOrder(order)
	}
	return result
}

// ToRepoOrders конвертирует срез моделей сервисного слоя в срез моделей репозитория
func ToRepoOrders(orders []*model.Order) []*repomodel.Order {
	if orders == nil {
		return nil
	}

	result := make([]*repomodel.Order, len(orders))
	for i, order := range orders {
		result[i] = ToRepoOrder(order)
	}
	return result
}
