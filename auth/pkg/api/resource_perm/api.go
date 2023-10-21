package resource_perm

import (
	resource_perm_create "auth/pkg/api/resource_perm/create"
	resource_perm_delete "auth/pkg/api/resource_perm/delete"
	resource_perm_get "auth/pkg/api/resource_perm/get"
	resource_perm_list "auth/pkg/api/resource_perm/list"
	resource_perm_update "auth/pkg/api/resource_perm/update"
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
	status, response := resource_perm_get.NewGetApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("resource_perm_id"),
		a.database,
	).Get()
	return ctx.Status(status).JSON(response)
}
func (a *Api) List(ctx *fiber.Ctx) (err error) {
	status, response := resource_perm_list.NewListApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		a.database,
	).List()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Create(ctx *fiber.Ctx) (err error) {
	status, response := resource_perm_create.NewCreateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Update(ctx *fiber.Ctx) (err error) {
	status, response := resource_perm_update.NewUpdateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Delete(ctx *fiber.Ctx) (err error) {
	status, response := resource_perm_delete.NewDeleteApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("resource_perm_id"),
		a.database,
	).Delete()
	return ctx.Status(status).JSON(response)
}
