package resource

import (
	resource_create "auth/pkg/api/resource/create"
	resource_delete "auth/pkg/api/resource/delete"
	resource_get "auth/pkg/api/resource/get"
	resource_list "auth/pkg/api/resource/list"
	resource_update "auth/pkg/api/resource/update"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Api struct {
	database *gorm.DB
}

func NewApi(
	database *gorm.DB,
) *Api {
	return &Api{
		database,
	}
}

func (a *Api) Get(ctx *fiber.Ctx) (err error) {
	status, response := resource_get.NewGetApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("resource_id"),
		a.database,
	).Get()
	return ctx.Status(status).JSON(response)
}
func (a *Api) List(ctx *fiber.Ctx) (err error) {
	status, response := resource_list.NewListApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		a.database,
	).List()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Create(ctx *fiber.Ctx) (err error) {
	status, response := resource_create.NewCreateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Update(ctx *fiber.Ctx) (err error) {
	status, response := resource_update.NewUpdateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Delete(ctx *fiber.Ctx) (err error) {
	status, response := resource_delete.NewDeleteApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("resource_id"),
		a.database,
	).Delete()
	return ctx.Status(status).JSON(response)
}
