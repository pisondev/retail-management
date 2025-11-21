package repository

import (
	"context"
	"database/sql"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	FindByID(ctx context.Context, tx *sql.Tx, userID ulid.ULID) (domain.User, error)
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.User, error)
	Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tx *sql.Tx, userID ulid.ULID) error
	AssignRole(ctx context.Context, tx *sql.Tx, userID ulid.ULID, roleName string) error
	UpdateRole(ctx context.Context, tx *sql.Tx, userID ulid.ULID, roleName string) error
}
