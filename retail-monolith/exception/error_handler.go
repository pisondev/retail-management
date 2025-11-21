package exception

import (
	"errors"
	"retail-management/model/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var validationErr validator.ValidationErrors

	code := fiber.StatusInternalServerError
	status := "INTERNAL SERVER ERROR"

	// 400 Bad Request
	if errors.As(err, &validationErr) {
		code = fiber.StatusBadRequest
		status = "BAD REQUEST"
	}
	if errors.Is(err, ErrInsufficientStock) {
		code = fiber.StatusBadRequest
		status = "BAD REQUEST"
	}

	// 401 Unauthorized
	if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrUnauthorizedLogin) {
		code = fiber.StatusUnauthorized
		status = "UNAUTHORIZED"
	}

	// 403 Forbidden
	if errors.Is(err, ErrForbidden) {
		code = fiber.StatusForbidden
		status = "FORBIDDEN"
	}

	// 404 Not Found
	if errors.Is(err, ErrNotFound) {
		code = fiber.StatusNotFound
		status = "NOT FOUND"
	}

	// 409 Conflict
	if errors.Is(err, ErrConflict) {
		code = fiber.StatusConflict
		status = "CONFLICT"
	}

	webResponse := web.WebResponse{
		Code:   code,
		Status: status,
		Data:   err.Error(),
	}

	return ctx.Status(code).JSON(webResponse)
}
