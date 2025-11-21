package controller

import "github.com/gofiber/fiber/v2"

type TransactionController interface {
	Create(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByID(ctx *fiber.Ctx) error
}
