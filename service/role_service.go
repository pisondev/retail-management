package service

import (
	"context"
	"retail-management/model/web"
)

type RoleService interface {
	FindAll(ctx context.Context) ([]web.RoleResponse, error)
}
