package repository

import (
	"context"
	"database/sql"
	"errors"
	"retail-management/model/domain"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct {
	Logger *logrus.Logger
}

func NewUserRepository(logger *logrus.Logger) UserRepository {
	return &UserRepositoryImpl{
		Logger: logger,
	}
}

func (repository *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	SQL := "INSERT INTO Users(user_id, username, hashed_password) VALUES (?, ?, ?)"

	repository.Logger.Info("---executing sql (register account)...")
	_, err := tx.ExecContext(ctx, SQL, user.UserID, user.Username, user.HashedPassword)
	if err != nil {
		repository.Logger.Errorf("---failed to execcontext: %v", err)
		return domain.User{}, err
	}
	repository.Logger.Info("---repository: success, returning back to service layer")
	return user, nil
}

func (repository *UserRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, userID ulid.ULID) (domain.User, error) {
	SQL := "SELECT user_id, username FROM Users WHERE user_id = ?"

	var user domain.User

	repository.Logger.Info("---executing sql (select by id)...")
	err := tx.QueryRowContext(ctx, SQL, userID).Scan(
		&user.UserID,
		&user.Username,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			repository.Logger.Warn("---cannot found user_id")
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
		}
		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	SQL := "SELECT user_id, username, hashed_password FROM Users WHERE username = ?"

	var user domain.User

	repository.Logger.Infof("---executing sql (select by username)...")
	err := tx.QueryRowContext(ctx, SQL, username).Scan(
		&user.UserID,
		&user.Username,
		&user.HashedPassword,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			repository.Logger.Warnf("---cannot found username: %v", username)
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
		}

		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.User, error) {
	SQL := "SELECT user_id, username FROM Users ORDER BY user_id ASC"

	repository.Logger.Info("---executing sql (select all)...")
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return []domain.User{}, err
	}
	defer rows.Close()

	repository.Logger.Info("---initializing users slice...")
	users := make([]domain.User, 0)
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.UserID,
			&user.Username,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.User{}, err
		}
		users = append(users, user)
	}
	repository.Logger.Infof("---users result: %v", users)
	repository.Logger.Info("---all rows has been scanned. returning back to service layer...")

	return users, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	SQL := "UPDATE Users SET username = ? WHERE user_id = ?"

	repository.Logger.Info("---executing sql (update user)...")
	_, err := tx.ExecContext(ctx, SQL, user.Username, user.UserID)
	if err != nil {
		repository.Logger.Errorf("---failed to update a user")
		return domain.User{}, err
	}

	repository.Logger.Info("successfully updated a user. trying to select...")
	SQLSelect := "SELECT user_id, username FROM Users WHERE user_id = ?"
	err = tx.QueryRowContext(ctx, SQLSelect, user.UserID).Scan(
		&user.UserID,
		&user.Username,
	)
	if err != nil {
		repository.Logger.Errorf("failed to scan a user: %v", err)
		return domain.User{}, err
	}

	repository.Logger.Info("successfully scan a user, returning back to service layer...")
	return user, nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, userID ulid.ULID) error {
	SQL := "DELETE FROM Users WHERE user_id = ?"

	repository.Logger.Info("---executing sql (delete)...")
	result, err := tx.ExecContext(ctx, SQL, userID)
	if err != nil {
		repository.Logger.Errorf("---failed to delete a user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.Logger.Errorf("---failed to check rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repository.Logger.Warnf("---failed to delete, user not found: %v", userID)
		return errors.New("cannot found user")
	}

	return nil
}
