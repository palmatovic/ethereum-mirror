package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-syncronizer/pkg/api/wallet_transaction/get"
	list "wallet-syncronizer/pkg/api/wallet_transaction/list"
	"wallet-syncronizer/pkg/util/url/wallet_transaction"
)

type Env struct {
	DB *gorm.DB
}

func NewEnv(db *gorm.DB) *Env {
	return &Env{DB: db}
}

func (e *Env) Get(ctx *fiber.Ctx) error {
	status, response := get.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(string(wallet_transaction.Id))).Get()
	return ctx.Status(status).JSON(response)
}

func (e *Env) List(ctx *fiber.Ctx) error {
	status, response := list.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB).List()
	return ctx.Status(status).JSON(response)
}
