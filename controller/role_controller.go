package controller

import "github.com/gofiber/fiber/v2"

type RoleController interface {
	FindAll(ctx *fiber.Ctx) error
}
