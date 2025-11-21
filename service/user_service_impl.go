package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"os"
	"retail-management/exception"
	"retail-management/helper"
	"retail-management/model/domain"
	"retail-management/model/web"
	"retail-management/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Logger         *logrus.Logger
}

func NewUserService(userRepository repository.UserRepository, db *sql.DB, validate *validator.Validate, logger *logrus.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             db,
		Validate:       validate,
		Logger:         logger,
	}
}

func (service *UserServiceImpl) Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserLoginResponse{}, err
	}

	service.Logger.Infof("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserLoginResponse{}, err
	}

	service.Logger.Infof("-executing repository.FindByUsername...")
	foundUser, err := service.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserLoginResponse{}, errRollback
		}
		return web.UserLoginResponse{}, err
	}

	service.Logger.Infof("-trying to hash the password...")
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(req.Password))
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserLoginResponse{}, errRollback
		}
		return web.UserLoginResponse{}, exception.ErrUnauthorizedLogin
	}

	service.Logger.Infof("-creating jwt claims...")
	claims := web.JWTClaims{
		UserID: foundUser.UserID.String(),
		Role:   foundUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   foundUser.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserLoginResponse{}, errRollback
		}
		service.Logger.Errorf("-failed to create tokenString")
		return web.UserLoginResponse{}, err
	}

	service.Logger.Infof("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		service.Logger.Errorf("-failed to commit tx: %v", errCommit)
		return web.UserLoginResponse{}, errCommit
	}

	service.Logger.Infof("-successfully commit tx!")
	return web.UserLoginResponse{
		Token: tokenString,
	}, nil
}

func (service *UserServiceImpl) FindByID(ctx context.Context, userID ulid.ULID) (web.UserResponse, error) {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}

	service.Logger.Info("-trying to execute r.FindByID...")
	foundUser, err := service.UserRepository.FindByID(ctx, tx, userID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserResponse{}, errRollback
		}
		if err == sql.ErrNoRows || errors.Is(err, exception.ErrNotFound) {
			service.Logger.Warn("-user with specific user_id doesn't exists")
			return web.UserResponse{}, exception.ErrNotFound
		}
		service.Logger.Errorf("-failed to execute r.FindByID: %v", err)
		return web.UserResponse{}, err
	}

	service.Logger.Infof("found a user with username: %v", foundUser.Username)

	return web.UserResponse{
		UserID:   foundUser.UserID,
		Username: foundUser.Username,
		Role:     foundUser.Role,
	}, nil
}

func (service *UserServiceImpl) Register(ctx context.Context, req web.UserAuthRequest) (web.UserRegisterResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserRegisterResponse{}, err
	}

	service.Logger.Infof("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserRegisterResponse{}, err
	}
	defer tx.Rollback()

	service.Logger.Infof("-executing repository.FindByUsername...")
	_, err = service.UserRepository.FindByUsername(ctx, tx, req.Username)

	if err != nil && !errors.Is(err, exception.ErrUnauthorizedLogin) && err != sql.ErrNoRows {
		service.Logger.Errorf("-failed to execute r.FindByUsername: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserRegisterResponse{}, errRollback
		}
		return web.UserRegisterResponse{}, err
	}
	if err == nil {
		service.Logger.Warnf("-username %v already exists.", req.Username)
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserRegisterResponse{}, errRollback
		}
		return web.UserRegisterResponse{}, exception.ErrConflict
	}

	service.Logger.Infof("-implementing ulid...")
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)

	service.Logger.Infof("-trying to hash the password...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		service.Logger.Errorf("-failed to hash the password")
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserRegisterResponse{}, errRollback
		}
		return web.UserRegisterResponse{}, err
	}

	user := domain.User{
		UserID:         ulid,
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
	}

	service.Logger.Infof("-executing repository.Save...")
	savedUser, err := service.UserRepository.Save(ctx, tx, user)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return web.UserRegisterResponse{}, errRollback
		}
		service.Logger.Errorf("failed to execute r.Save: %v", err)
		return web.UserRegisterResponse{}, err
	}

	targetRole := "cashier"
	if req.Role != "" {
		targetRole = req.Role
	}

	service.Logger.Infof("-executing repository.AssignRole (%s)...", targetRole)
	err = service.UserRepository.AssignRole(ctx, tx, savedUser.UserID, targetRole)
	if err != nil {
		service.Logger.Errorf("failed to execute r.AssignRole: %v", err)
		return web.UserRegisterResponse{}, err
	}

	service.Logger.Infof("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		service.Logger.Errorf("failed to commit tx: %v", errCommit)
		return web.UserRegisterResponse{}, errCommit
	}

	service.Logger.Infof("-successfully commit tx!")

	savedUserResponse := helper.ToUserRegisterResponse(savedUser)
	savedUserResponse.Role = targetRole

	return savedUserResponse, nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]web.UserResponse, error) {
	service.Logger.Info("trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return []web.UserResponse{}, err
	}

	foundUsers, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		service.Logger.Errorf("failed to execute r.FindAll: %v", err)
		errRollback := tx.Rollback()
		if errRollback != nil {
			return []web.UserResponse{}, err
		}
		return []web.UserResponse{}, err
	}

	var responses []web.UserResponse
	for _, u := range foundUsers {
		responses = append(responses, web.UserResponse{
			UserID:   u.UserID,
			Username: u.Username,
			Role:     u.Role,
		})
	}

	return responses, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, req web.UserUpdateRequest) (web.UserResponse, error) {
	err := service.Validate.Struct(req)
	if err != nil {
		return web.UserResponse{}, err
	}

	service.Logger.Infof("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return web.UserResponse{}, err
	}
	defer tx.Rollback()

	selectedUser, err := service.UserRepository.FindByID(ctx, tx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows || errors.Is(err, exception.ErrNotFound) {
			return web.UserResponse{}, exception.ErrNotFound
		}
		return web.UserResponse{}, err
	}

	if req.Username != nil {
		selectedUser.Username = *req.Username
		_, err := service.UserRepository.Update(ctx, tx, selectedUser)
		if err != nil {
			service.Logger.Errorf("-failed to update user info: %v", err)
			return web.UserResponse{}, err
		}
	}

	if req.Role != nil {
		service.Logger.Infof("-updating role to %s...", *req.Role)
		err := service.UserRepository.UpdateRole(ctx, tx, selectedUser.UserID, *req.Role)
		if err != nil {
			service.Logger.Errorf("-failed to update user role: %v", err)
			return web.UserResponse{}, err
		}
	}

	service.Logger.Info("-trying to commit tx...")
	if err = tx.Commit(); err != nil {
		return web.UserResponse{}, err
	}

	finalRole := selectedUser.Role
	if req.Role != nil {
		finalRole = *req.Role
	}

	return web.UserResponse{
		UserID:   selectedUser.UserID,
		Username: selectedUser.Username,
		Role:     finalRole,
	}, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, userID ulid.ULID) error {
	service.Logger.Info("-trying to begin tx...")
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}

	service.Logger.Info("-trying to execute r.Delete...")
	err = service.UserRepository.Delete(ctx, tx, userID)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}

		if errors.Is(err, exception.ErrNotFound) {
			service.Logger.Warnf("-user not found for delete: %v", userID)
			return exception.ErrNotFound
		}

		service.Logger.Errorf("-failed to execute r.Delete: %v", err)
		return err
	}

	service.Logger.Info("-trying to commit tx...")
	errCommit := tx.Commit()
	if errCommit != nil {
		return errCommit
	}

	service.Logger.Info("-successfully delete the specific user")
	return nil
}
