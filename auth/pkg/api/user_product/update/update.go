package update

import (
	user_product_model "auth/pkg/model/api/user_product/update"
	util_json "auth/pkg/model/json"
	user_product_update_service "auth/pkg/service/user_product/update"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Api struct {
	db     *gorm.DB
	body   []byte
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, body []byte) *Api {
	return &Api{
		body:   body,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "body": string(body)},
	}
}

func (a *Api) Update() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var userProduct user_product_model.UserProduct
	if err := json.Unmarshal(a.body, &userProduct); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, groupDb, err := user_product_update_service.NewService(a.db, &userProduct).Update()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"user_product": groupDb}}
}
