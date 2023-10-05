package get

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"wallet-synchronizer/pkg/api/wallet/create"
	"wallet-synchronizer/pkg/api/wallet/delete"
	"wallet-synchronizer/pkg/api/wallet/get"
	list "wallet-synchronizer/pkg/api/wallet/list"
	"wallet-synchronizer/pkg/api/wallet/update"
	"wallet-synchronizer/pkg/util/url/wallet"
)

type Env struct {
	DB *gorm.DB
}

func NewEnv(db *gorm.DB) *Env {
	return &Env{DB: db}
}

func (e *Env) Get(ctx *fiber.Ctx) error {
	status, response := get.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(string(wallet.Id))).Get()
	return ctx.Status(status).JSON(response)
}
func (e *Env) List(ctx *fiber.Ctx) error {
	status, response := list.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB).List()
	return ctx.Status(status).JSON(response)
}

func (e *Env) Create(ctx *fiber.Ctx) error {
	status, response := create.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Body()).Create()
	return ctx.Status(status).JSON(response)
}

func (e *Env) Update(ctx *fiber.Ctx) error {
	status, response := update.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Body()).Update()
	return ctx.Status(status).JSON(response)
}

func (e *Env) Delete(ctx *fiber.Ctx) error {
	status, response := delete.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(string(wallet.Id))).Delete()
	return ctx.Status(status).JSON(response)
}
