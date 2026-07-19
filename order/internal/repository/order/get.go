package order

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

// FindByUUID находит заказ по UUID в PostgreSQL
func (r *PostgresRepository) FindByUUID(ctx context.Context, uuid string) (*model.Order, error) {
	query := `
		SELECT 
			order_uuid, user_uuid, part_uuids, total_price,
			transaction_uuid, payment_method, status, created_at, updated_at
		FROM orders
		WHERE order_uuid = $1
	`

	var order model.Order
	var partUUIDsJSON []byte
	var transactionUUID *string
	var paymentMethod *string

	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&partUUIDsJSON,
		&order.TotalPrice,
		&transactionUUID,
		&paymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}

	// Распарсим JSONB в []string
	if err := json.Unmarshal(partUUIDsJSON, &order.PartUUIDs); err != nil {
		return nil, err
	}

	if transactionUUID != nil {
		order.TransactionUUID = transactionUUID
	}
	if paymentMethod != nil {
		method := model.PaymentMethod(*paymentMethod)
		order.PaymentMethod = &method
	}

	return &order, nil
}