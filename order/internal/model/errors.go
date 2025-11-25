package model

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderUUID  = errors.New("invalid order UUID")
	ErrOrderAlreadyPaid  = errors.New("order already paid")
	ErrOrderCancelled    = errors.New("order cancelled")
	ErrOrderFulfilled    = errors.New("order fulfilled")
	ErrPaymentRequired   = errors.New("payment method required")
	ErrPartsNotSpecified = errors.New("at least one part must be specified")
	ErrPartsNotFound     = errors.New("some parts not found")
)
