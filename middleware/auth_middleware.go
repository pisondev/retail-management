package middleware

import (
	"os"
	"retail-management/model/web"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			webResponse := web.WebResponse{
				Code:   fiber.StatusUnauthorized,
				Status: "UNAUTHORIZED",
				Data:   "Missing JWT",
			}
			return ctx.Status(fiber.StatusUnauthorized).JSON(webResponse)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &web.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			webResponse := web.WebResponse{
				Code:   fiber.StatusUnauthorized,
				Status: "UNAUTHORIZED",
				Data:   "Invalid or expired JWT",
			}
			return ctx.Status(fiber.StatusUnauthorized).JSON(webResponse)
		}

		ctx.Locals("userID", claims.UserID)
		ctx.Locals("role", claims.Role)

		return ctx.Next()
	}
}
