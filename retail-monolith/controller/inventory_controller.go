package controller

import "github.com/gofiber/fiber/v2"

type InventoryLogController interface {
	Adjust(ctx *fiber.Ctx) error
}
