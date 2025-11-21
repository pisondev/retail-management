package controller

import "github.com/gofiber/fiber/v2"

type UserController interface {
	// Auth
	Login(ctx *fiber.Ctx) error
	GetMe(ctx *fiber.Ctx) error

	// user management
	Register(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
