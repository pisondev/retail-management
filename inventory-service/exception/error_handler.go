package exception

import (
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCErrorHandler(logger *logrus.Logger, err error, contextMessage string) error {
	if err == nil {
		return nil
	}

	logger.Errorf("%s: %v", contextMessage, err)

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound) {
		return status.Error(codes.NotFound, ErrNotFound.Error())
	}

	if errors.Is(err, ErrInvalidID) || errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrStockNegative) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Internal, ErrInternalServer.Error())
}

func FormatErrorMessage(logger *logrus.Logger, err error, contextMessage string) string {
	logger.Errorf("%s: %v", contextMessage, err)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound.Error()
	}

	if errors.Is(err, ErrInsufficientStock) || errors.Is(err, ErrStockNegative) || errors.Is(err, ErrInvalidID) {
		return err.Error()
	}

	return ErrInternalServer.Error()
}
