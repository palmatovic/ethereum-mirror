package user_resource_perm_perm

import (
	user_resource_perm_create "auth/pkg/api/user_resource_perm/create"
	user_resource_perm_delete "auth/pkg/api/user_resource_perm/delete"
	user_resource_perm_get "auth/pkg/api/user_resource_perm/get"
	user_resource_perm_list "auth/pkg/api/user_resource_perm/list"
	user_resource_perm_update "auth/pkg/api/user_resource_perm/update"
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
	status, response := user_resource_perm_get.NewGetApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("user_resource_perm_id"),
		a.database,
	).Get()
	return ctx.Status(status).JSON(response)
}
func (a *Api) List(ctx *fiber.Ctx) (err error) {
	status, response := user_resource_perm_list.NewListApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		a.database,
	).List()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Create(ctx *fiber.Ctx) (err error) {
	status, response := user_resource_perm_create.NewCreateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Update(ctx *fiber.Ctx) (err error) {
	status, response := user_resource_perm_update.NewUpdateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Delete(ctx *fiber.Ctx) (err error) {
	status, response := user_resource_perm_delete.NewDeleteApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("user_resource_perm_id"),
		a.database,
	).Delete()
	return ctx.Status(status).JSON(response)
}
