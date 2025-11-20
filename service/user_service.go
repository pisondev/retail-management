package service

import (
	"context"
	"retail-management/model/web"

	"github.com/oklog/ulid/v2"
)

type UserService interface {
	Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error)
	FindByID(ctx context.Context, userID ulid.ULID) (web.UserResponse, error)
	Register(ctx context.Context, req web.UserAuthRequest) (web.UserRegisterResponse, error)
	FindAll(ctx context.Context) ([]web.UserResponse, error)
	Update(ctx context.Context, req web.UserUpdateRequest) (web.UserResponse, error)
	Delete(ctx context.Context, userID ulid.ULID) error
}
