package update

import (
	model_update_product "auth/pkg/model/api/product/update"
	util_json "auth/pkg/model/json"
	product_update_service "auth/pkg/service/product/update"
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
	var product model_update_product.Product
	if err := json.Unmarshal(a.body, &product); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, productDb, err := product_update_service.NewService(a.db, &product).Update()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"product": productDb}}
}
