package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_token_get_service "wallet-syncronizer/pkg/service/wallet_token/get"
	"wallet-syncronizer/pkg/util/json"
)

type Api struct {
	db       *gorm.DB
	walletId string
	tokenId  string
	fields   logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, walletId string, tokenId string) *Api {
	return &Api{
		walletId: walletId,
		db:       db,
		fields:   logrus.Fields{"uuid": uuid, "url": url, "wallet_id": walletId, "token_id": tokenId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.walletId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty wallet_id")
	}
	if len(a.tokenId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id")
	}
	httpStatus, walletToken, err := wallet_token_get_service.NewService(a.db, a.walletId, a.tokenId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"wallet_token": walletToken}}
}
