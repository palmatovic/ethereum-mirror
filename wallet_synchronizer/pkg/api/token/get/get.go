package get

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wallet-synchronizer/pkg/model/json"
	token_get_service "wallet-synchronizer/pkg/service/token/get"
)

type Api struct {
	db      *gorm.DB
	tokenId string
	fields  logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, tokenId string) *Api {
	return &Api{
		tokenId: tokenId,
		db:      db,
		fields:  logrus.Fields{"uuid": uuid, "url": url, "token_id": tokenId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.tokenId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty token_id")).Errorf("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id")
	}
	httpStatus, token, err := token_get_service.NewService(a.db, a.tokenId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"token": token}}
}
