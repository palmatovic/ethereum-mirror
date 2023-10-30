package user

import (
	"auth/pkg/api/user/create"
	"auth/pkg/api/user/delete"
	"auth/pkg/api/user/get"
	list "auth/pkg/api/user/list"
	"auth/pkg/api/user/update"
	"auth/pkg/url/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Api struct {
	DB *gorm.DB
}

func NewApi(db *gorm.DB) *Api {
	return &Api{DB: db}
}

func (e *Api) Get(ctx *fiber.Ctx) error {
	status, response := get.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(user.Id)).Get()
	return ctx.Status(status).JSON(response)
}
func (e *Api) List(ctx *fiber.Ctx) error {
	status, response := list.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Query("page_size"), ctx.Query("page_number")).List()
	return ctx.Status(status).JSON(response)
}

func (e *Api) Create(ctx *fiber.Ctx) error {
	status, response := create.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Body()).Create()
	return ctx.Status(status).JSON(response)
}

func (e *Api) Update(ctx *fiber.Ctx) error {
	status, response := update.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Body()).Update()
	return ctx.Status(status).JSON(response)
}

func (e *Api) Delete(ctx *fiber.Ctx) error {
	status, response := delete.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), e.DB, ctx.Params(user.Id)).Delete()
	return ctx.Status(status).JSON(response)
}
