package json

import "github.com/gofiber/fiber/v2"

type Response struct {
	Error *fiber.Error `json:"error,omitempty"`
	Data  interface{}  `json:"data,omitempty"`
}

func NewErrorResponse(httpStatus int, message string) Response {
	return Response{
		Error: fiber.NewError(httpStatus, message),
	}
}
