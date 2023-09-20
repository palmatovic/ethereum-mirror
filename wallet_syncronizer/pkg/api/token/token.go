package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-syncronizer/pkg/api/token/get"
	list "wallet-syncronizer/pkg/api/token/list"
	"wallet-syncronizer/pkg/util/url/token"
)

type Env struct {
	DB *gorm.DB
}

func NewEnv(db *gorm.DB) *Env {
	return &Env{DB: db}
}

func (e *Env) Get(ctx *fiber.Ctx) error {
	status, response := get.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(string(token.Id))).Get()
	return ctx.Status(status).JSON(response)
}

func (e *Env) GetList(ctx *fiber.Ctx) error {
	return list.GetList(ctx, e.DB)
}
