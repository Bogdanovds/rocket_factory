package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bogdanovds/rocket_factory/order/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, orderID uuid.UUID, paymentMethod string) (*model.Order, error) {
	order, err := s.repo.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return nil, model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return nil, model.ErrOrderCancelled
	case model.OrderStatusFulfilled:
		return nil, model.ErrOrderFulfilled
	}

	if paymentMethod == "" {
		return nil, model.ErrPaymentRequired
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, orderID, order.UserID, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = paymentMethod
	order.TransactionID = transactionID

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return order, nil
}
