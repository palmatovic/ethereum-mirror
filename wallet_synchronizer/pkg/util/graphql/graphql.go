package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/gqlerrors"
)

func MapGraphQLErrorToHTTPStatus(gqlError *gqlerrors.FormattedError) int {
	switch gqlError.Message {
	case "record not found":
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}
