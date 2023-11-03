package validator

import (
	"crypto/rsa"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	util_json "wallet-synchronizer/pkg/model/json"
)

type Service struct {
	jwtPublicKey *rsa.PublicKey
}

func NewService(jwtPublicKey *rsa.PublicKey) *Service {
	return &Service{
		jwtPublicKey: jwtPublicKey,
	}
}

func (s *Service) Validator() fiber.Handler {
	jwtValidator := jwtware.New(jwtware.Config{
		SigningMethod: "RS256",
		SigningKey:    s.jwtPublicKey,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			uuidReq := ctx.Locals("uuid").(string)
			url := ctx.OriginalURL()
			var fields logrus.Fields
			fields["url"] = url
			fields["uuid"] = uuidReq
			for k, v := range ctx.GetReqHeaders() {
				fields[k] = v
			}
			logrus.WithFields(fields).WithError(err).Errorf("terminated with failure")
			if errors.Is(err, jwt.ErrTokenExpired) {
				return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, "expired token"))
			}
			return ctx.Status(fiber.StatusUnauthorized).JSON(util_json.NewErrorResponse(fiber.StatusUnauthorized, "invalid token"))
		},
	})
	return jwtValidator
}
