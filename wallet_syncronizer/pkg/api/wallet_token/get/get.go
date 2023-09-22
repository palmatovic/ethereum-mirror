package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_token_get_service "wallet-syncronizer/pkg/service/wallet_token/get"
	"wallet-syncronizer/pkg/util/json"
)

type Api struct {
	db            *gorm.DB
	walletTokenId string
	fields        logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, walletTokenId string) *Api {
	return &Api{
		walletTokenId: walletTokenId,
		db:            db,
		fields:        logrus.Fields{"uuid": uuid, "url": url, "wallet_token_id": walletTokenId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.walletTokenId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id")
	}
	httpStatus, walletToken, err := wallet_token_get_service.NewService(a.db, a.walletTokenId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"wallet_token": walletToken}}
}
