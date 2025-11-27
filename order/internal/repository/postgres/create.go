package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// Create создаёт новый заказ в базе данных
func (r *Repository) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (id, user_id, part_ids, total_price, status, payment_method, transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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

	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.UserID,
		pq.Array(partIDs),
		order.TotalPrice,
		string(order.Status),
		order.PaymentMethod,
		transactionID,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}
