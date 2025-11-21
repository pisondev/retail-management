package middleware

import (
	"retail-management/exception"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		roleRaw := ctx.Locals("role")
		if roleRaw == nil {
			return exception.ErrUnauthorized
		}

		role, ok := roleRaw.(string)
		if !ok {
			return exception.ErrUnauthorized
		}

		if role != "admin" {
			return exception.ErrForbidden
		}

		return ctx.Next()
	}
}
