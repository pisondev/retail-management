package repository

import (
	"context"
	"database/sql"
	"errors"
	"retail-management/exception"
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
	SQL := `
        SELECT 
            u.user_id, 
            u.username, 
            r.role_name
        FROM Users u
        JOIN User_Roles ur ON u.user_id = ur.user_id
        JOIN Roles r ON ur.role_id = r.role_id
        WHERE u.user_id = ?
    `

	var user domain.User

	repository.Logger.Info("---executing sql (select by id with role)...")
	err := tx.QueryRowContext(ctx, SQL, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			repository.Logger.Warnf("---user not found by id: %v", userID)
			return domain.User{}, exception.ErrNotFound
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return domain.User{}, err
		}
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	SQL := `
        SELECT 
            u.user_id, 
            u.username, 
            u.hashed_password, 
            r.role_name
        FROM Users u
        JOIN User_Roles ur ON u.user_id = ur.user_id
        JOIN Roles r ON ur.role_id = r.role_id
        WHERE u.username = ?
    `
	var user domain.User

	repository.Logger.Infof("---executing sql (select by username)...")
	err := tx.QueryRowContext(ctx, SQL, username).Scan(
		&user.UserID,
		&user.Username,
		&user.HashedPassword,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			repository.Logger.Warnf("---cannot found username: %v", username)
			return domain.User{}, exception.ErrUnauthorizedLogin
		} else {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return domain.User{}, err
		}
	}

	return user, nil
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.User, error) {
	SQL := `
        SELECT 
            u.user_id, 
            u.username,
            r.role_name
        FROM Users u
        JOIN User_Roles ur ON u.user_id = ur.user_id
        JOIN Roles r ON ur.role_id = r.role_id
        ORDER BY u.username ASC
    `

	repository.Logger.Info("---executing sql (select all users with roles)...")
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
			&user.Role,
		)
		if err != nil {
			repository.Logger.Errorf("---failed to scan row: %v", err)
			return []domain.User{}, err
		}
		users = append(users, user)
	}

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

	return repository.FindByID(ctx, tx, user.UserID)
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
		return exception.ErrNotFound
	}

	return nil
}

func (repository *UserRepositoryImpl) AssignRole(ctx context.Context, tx *sql.Tx, userID ulid.ULID, roleName string) error {
	SQL := `INSERT INTO User_Roles (user_id, role_id) 
            SELECT ?, role_id FROM Roles WHERE role_name = ?`

	repository.Logger.Infof("---executing sql (assign role '%s')...", roleName)
	result, err := tx.ExecContext(ctx, SQL, userID, roleName)
	if err != nil {
		repository.Logger.Errorf("---failed to assign role: %v", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("role not found")
	}

	return nil
}

func (repository *UserRepositoryImpl) UpdateRole(ctx context.Context, tx *sql.Tx, userID ulid.ULID, roleName string) error {
	SQL := `UPDATE User_Roles 
            SET role_id = (SELECT role_id FROM Roles WHERE role_name = ?) 
            WHERE user_id = ?`

	repository.Logger.Infof("---executing sql (update role to '%s')...", roleName)
	result, err := tx.ExecContext(ctx, SQL, roleName, userID)
	if err != nil {
		repository.Logger.Errorf("---failed to update user role: %v", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return repository.AssignRole(ctx, tx, userID, roleName)
	}

	return nil
}
