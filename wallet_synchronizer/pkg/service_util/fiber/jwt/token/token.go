package token

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	util_json "wallet-synchronizer/pkg/model/json"
	util_jwt "wallet-synchronizer/pkg/service_util/jwt"
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

func HasResourcePerm(resource string, perm string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
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
		var idxProduct int
		if idxProduct = slices.IndexFunc(token.Products, func(product util_jwt.Product) bool {
			return product.Name == "eth_mirror"
		}); idxProduct == -1 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, "auth product not found"))
		}
		for _, group := range token.Products[idxProduct].Groups {
			for _, role := range group.Roles {
				var idxResource int
				if idxResource = slices.IndexFunc(role.Resources, func(res util_jwt.Resource) bool {
					return res.Name == resource
				}); idxResource == -1 {
					return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, fmt.Sprintf("%s resource not found", resource)))
				}
				var idxPerm int
				if idxPerm = slices.IndexFunc(role.Resources[idxResource].Perms, func(p util_jwt.Perm) bool {
					return p.Id == perm
				}); idxPerm == -1 {
					return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, fmt.Sprintf("%s perm not found", perm)))
				}
			}
		}
		return ctx.Next()
	}
}
