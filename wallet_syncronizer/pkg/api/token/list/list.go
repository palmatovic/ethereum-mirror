package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	token_list_service "wallet-syncronizer/pkg/service/token/list"
	"wallet-syncronizer/pkg/util/json"
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
