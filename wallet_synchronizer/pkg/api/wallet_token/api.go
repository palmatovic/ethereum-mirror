package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-synchronizer/pkg/api/wallet_token/get"
	list "wallet-synchronizer/pkg/api/wallet_token/list"
	token_url "wallet-synchronizer/pkg/url/token"
	wallet_url "wallet-synchronizer/pkg/url/wallet"
)

type Api struct {
	DB *gorm.DB
}

func NewApi(db *gorm.DB) *Api {
	return &Api{DB: db}
}

func (e *Api) Get(ctx *fiber.Ctx) error {
	status, response := get.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(string(wallet_url.Id)), ctx.Params(string(token_url.Id))).Get()
	return ctx.Status(status).JSON(response)
}

func (e *Api) List(ctx *fiber.Ctx) error {
	status, response := list.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB).List()
	return ctx.Status(status).JSON(response)
}
