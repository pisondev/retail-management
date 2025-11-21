package service

import (
	"context"
	"retail-management/model/web"
)

type InventoryLogService interface {
	Adjust(ctx context.Context, req web.InventoryLogRequest) (web.InventoryLogResponse, error)
}
