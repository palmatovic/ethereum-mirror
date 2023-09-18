package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-syncronizer/pkg/util/json"
)

func Get(ctx *fiber.Ctx, db *gorm.DB) error {
	var tokenId string
	if tokenId = ctx.Params("token_id"); len(tokenId) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id"))
	}
	return ctx.Status(200).JSON(nil)
}
