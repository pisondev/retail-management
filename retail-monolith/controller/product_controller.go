package controller

import "github.com/gofiber/fiber/v2"

type ProductController interface {
	Create(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByID(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	UpdateStock(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
