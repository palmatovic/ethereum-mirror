package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wallet-synchronizer/pkg/model/json"
	token_list_service "wallet-synchronizer/pkg/service/token/list"
)

type Api struct {
	db     *gorm.DB
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB) *Api {
	return &Api{
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url},
	}
}

func (a *Api) List() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	httpStatus, token, err := token_list_service.NewService(a.db).List()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"token": token}}
}
