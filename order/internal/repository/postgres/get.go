package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

// Get получает заказ по ID из базы данных
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	query := `
		SELECT id, user_id, part_ids, total_price, status, payment_method, transaction_id
		FROM orders
		WHERE id = $1
	`

	var order model.Order
	var partIDs pq.StringArray
	var paymentMethod sql.NullString
	var transactionID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&partIDs,
		&order.TotalPrice,
		&order.Status,
		&paymentMethod,
		&transactionID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Конвертируем []string в []uuid.UUID
	order.PartIDs = make([]uuid.UUID, len(partIDs))
	for i, idStr := range partIDs {
		parsedID, parseErr := uuid.Parse(idStr)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse part ID: %w", parseErr)
		}
		order.PartIDs[i] = parsedID
	}

	if paymentMethod.Valid {
		order.PaymentMethod = paymentMethod.String
	}

	if transactionID.Valid {
		parsedTransactionID, parseErr := uuid.Parse(transactionID.String)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse transaction ID: %w", parseErr)
		}
		order.TransactionID = parsedTransactionID
	}

	return &order, nil
}
