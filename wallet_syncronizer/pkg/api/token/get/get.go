package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	token_service "wallet-syncronizer/pkg/service/token"
	"wallet-syncronizer/pkg/util/json"
)

type Api struct {
	db      *gorm.DB
	uuid    string
	url     string
	tokenId string
	fields  logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, tokenId string) *Api {
	return &Api{
		tokenId: tokenId,
		db:      db,
		uuid:    uuid,
		url:     url,
		fields:  logrus.Fields{"uuid": uuid, "url": url, "token_id": tokenId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.tokenId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id")
	}
	httpStatus, token, err := token_service.GetToken(a.db, a.tokenId)
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"token": token}}
}
