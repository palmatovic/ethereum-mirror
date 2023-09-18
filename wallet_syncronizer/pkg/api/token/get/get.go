package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Get(ctx *fiber.Ctx, db *gorm.DB) error {
	return ctx.Status(200).JSON(nil)
}
