package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// Update обновляет заказ в базе данных
func (r *Repository) Update(ctx context.Context, order *model.Order) error {
	query := `
		UPDATE orders
		SET user_id = $2, part_ids = $3, total_price = $4, status = $5, 
		    payment_method = $6, transaction_id = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	// Конвертируем []uuid.UUID в []string для pq.Array
	partIDs := make([]string, len(order.PartIDs))
	for i, id := range order.PartIDs {
		partIDs[i] = id.String()
	}

	var transactionID interface{}
	if order.TransactionID.String() != "00000000-0000-0000-0000-000000000000" {
		transactionID = order.TransactionID
	}

	result, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		pq.Array(partIDs),
		order.TotalPrice,
		string(order.Status),
		order.PaymentMethod,
		transactionID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}

