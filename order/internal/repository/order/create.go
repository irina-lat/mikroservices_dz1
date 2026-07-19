package order

import (
	"context"
	"encoding/json"
	"time"

	"github.com/irina-lat/microservices-course/order/internal/model"
)

func (r *PostgresRepository) Save(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (
			order_uuid, user_uuid, part_uuids, total_price,
			transaction_uuid, payment_method, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	partUUIDsJSON, err := json.Marshal(order.PartUUIDs)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = r.db.Exec(ctx, query,
		order.OrderUUID,
		order.UserUUID,
		partUUIDsJSON,
		order.TotalPrice,
		order.TransactionUUID,
		order.PaymentMethod,
		order.Status,
		now,
		now,
	)
	return err
}
