package model

import "errors"

// Ошибки сервиса оплаты
var (
	ErrInvalidOrderUUID     = errors.New("invalid order UUID")
	ErrInvalidUserUUID      = errors.New("invalid user UUID")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPaymentFailed        = errors.New("payment processing failed")
	ErrTransactionNotFound  = errors.New("transaction not found")
)
