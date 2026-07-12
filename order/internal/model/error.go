package model

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyPaid    = errors.New("order already paid")
	ErrOrderAlreadyCanceled = errors.New("order already canceled")
	ErrPartNotFound        = errors.New("some parts not found")
)