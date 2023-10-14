package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/gqlerrors"
	"strings"
)

func MapGraphQLErrorToHTTPStatus(gqlError *gqlerrors.FormattedError) int {
	switch gqlError.Message {
	case "record not found":
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

func Stringify(jsonString string) interface{} {
	cleaned := strings.ReplaceAll(jsonString, "\t", " ")
	return strings.ReplaceAll(cleaned, "\n", "")
}
