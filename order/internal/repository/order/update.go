package order

import (
	"context"
	"encoding/json"
	"time"

	"order/internal/model"
)

// Update обновляет заказ в PostgreSQL
func (r *PostgresRepository) Update(ctx context.Context, order *model.Order) error {
	query := `
		UPDATE orders 
		SET 
			user_uuid = $1,
			part_uuids = $2,
			total_price = $3,
			transaction_uuid = $4,
			payment_method = $5,
			status = $6,
			updated_at = $7
		WHERE order_uuid = $8
	`

	partUUIDsJSON, err := json.Marshal(order.PartUUIDs)
	if err != nil {
		return err
	}

	var transactionUUID interface{} = nil
	if order.TransactionUUID != nil {
		transactionUUID = *order.TransactionUUID
	}

	var paymentMethod interface{} = nil
	if order.PaymentMethod != nil {
		paymentMethod = string(*order.PaymentMethod)
	}

	result, err := r.db.Exec(ctx, query,
		order.UserUUID,
		partUUIDsJSON,
		order.TotalPrice,
		transactionUUID,
		paymentMethod,
		string(order.Status),
		time.Now(),
		order.OrderUUID,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}
	return nil
}
