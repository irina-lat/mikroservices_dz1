package model

import "errors"

var (
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrEmptyOrderUUID       = errors.New("order_uuid is required")
	ErrEmptyUserUUID        = errors.New("user_uuid is required")
)