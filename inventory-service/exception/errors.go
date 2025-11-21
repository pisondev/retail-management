package exception

import "errors"

var (
	ErrInvalidID         = errors.New("invalid id format")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrStockNegative     = errors.New("stock cannot be negative")
	ErrInsufficientStock = errors.New("insufficient stock")

	ErrNotFound = errors.New("data not found")

	ErrInternalServer = errors.New("internal server error")
	ErrDatabase       = errors.New("database operation failed")
)
