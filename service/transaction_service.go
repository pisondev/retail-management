package service

import (
	"context"
	"retail-management/model/web"

	"github.com/oklog/ulid/v2"
)

type TransactionService interface {
	Create(ctx context.Context, req web.TransactionRequest) (web.TransactionResponse, error)
	FindAll(ctx context.Context, requesterUserID ulid.ULID, requesterRole string) ([]web.TransactionResponse, error)
	FindByID(ctx context.Context, requesterUserID ulid.ULID, requesterRole string, transactionID ulid.ULID) (web.TransactionResponse, error)
}
