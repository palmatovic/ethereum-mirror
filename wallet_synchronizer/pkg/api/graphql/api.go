package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"wallet-synchronizer/pkg/api/graphql/post"
)

type Api struct {
	schema graphql.Schema
}

func NewApi(schema graphql.Schema) *Api {
	return &Api{schema: schema}
}

func (a *Api) Post(ctx *fiber.Ctx) error {
	status, response := post.NewApi(ctx.Locals("uuid").(string), ctx.OriginalURL(), a.schema, string(ctx.Body())).Post()
	return ctx.Status(status).JSON(response)
}
