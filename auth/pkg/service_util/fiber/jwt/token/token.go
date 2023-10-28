package token

import (
	util_json "auth/pkg/model/json"
	util_jwt "auth/pkg/service_util/jwt"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func AccessTokenValidator(ctx *fiber.Ctx) error {
	uuidReq := ctx.Locals("uuid").(string)
	url := ctx.OriginalURL()
	var fields logrus.Fields
	fields["url"] = url
	fields["uuid"] = uuidReq
	for k, v := range ctx.GetReqHeaders() {
		fields[k] = v
	}
	token, err := util_jwt.NewService(ctx).Extract()
	if err != nil {
		logrus.WithFields(fields).WithError(err).Errorf("terminated with failure")
		return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, err.Error()))
	}
	if !token.A {
		logrus.WithFields(fields).WithError(errors.New("not access token")).Errorf("terminated with error")
		return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, "invalid token"))
	}
	return ctx.Next()
}

func RefreshTokenValidator(ctx *fiber.Ctx) error {
	uuidReq := ctx.Locals("uuid").(string)
	url := ctx.OriginalURL()
	var fields logrus.Fields
	fields["url"] = url
	fields["uuid"] = uuidReq
	for k, v := range ctx.GetReqHeaders() {
		fields[k] = v
	}
	token, err := util_jwt.NewService(ctx).Extract()
	if err != nil {
		logrus.WithFields(fields).WithError(err).Errorf("terminated with failure")
		return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, err.Error()))
	}
	if token.A {
		logrus.WithFields(fields).WithError(errors.New("not refresh token")).Errorf("terminated with error")
		return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, "invalid token"))
	}
	return ctx.Next()
}
