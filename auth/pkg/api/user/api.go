package user

import (
	user_create "auth/pkg/api/user/create"
	user_delete "auth/pkg/api/user/delete"
	user_get "auth/pkg/api/user/get"
	user_list "auth/pkg/api/user/list"
	user_update "auth/pkg/api/user/update"
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
	status, response := user_get.NewGetApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("user_id"),
		a.database,
	).Get()
	return ctx.Status(status).JSON(response)
}
func (a *Api) List(ctx *fiber.Ctx) (err error) {
	status, response := user_list.NewListApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		a.database,
	).List()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Create(ctx *fiber.Ctx) (err error) {
	status, response := user_create.NewCreateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Update(ctx *fiber.Ctx) (err error) {
	status, response := user_update.NewUpdateApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Body(),
		a.database,
	).Update()
	return ctx.Status(status).JSON(response)
}
func (a *Api) Delete(ctx *fiber.Ctx) (err error) {
	status, response := user_delete.NewDeleteApi(
		ctx.OriginalURL(),
		ctx.Locals("uuid").(string),
		token,
		ctx.Params("user_id"),
		a.database,
	).Delete()
	return ctx.Status(status).JSON(response)
}
