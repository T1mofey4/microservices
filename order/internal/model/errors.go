package model

import "errors"

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderAlreadyPaid     = errors.New("order is already paid")
	ErrOrderCancelled       = errors.New("order cancelled")
	ErrPartNotFound         = errors.New("part not found")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
