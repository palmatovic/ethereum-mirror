package create

import (
	model_create_product "auth/pkg/model/api/product/create"
	util_json "auth/pkg/model/json"
	product_create_service "auth/pkg/service/product/create"
	"auth/pkg/service_util/aes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Api struct {
	db                  *gorm.DB
	body                []byte
	fields              logrus.Fields
	aes256EncryptionKey *aes.Key
}

func NewApi(uuid string, url string, db *gorm.DB, aes256EncryptionKey *aes.Key, body []byte) *Api {
	return &Api{
		body:                body,
		db:                  db,
		fields:              logrus.Fields{"uuid": uuid, "url": url, "body": util_json.Stringify(body)},
		aes256EncryptionKey: aes256EncryptionKey,
	}
}

func (a *Api) Create() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var product model_create_product.Product
	if err := json.Unmarshal(a.body, &product); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, productDb, err := product_create_service.NewService(a.db, &product, a.aes256EncryptionKey).Create()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"product": productDb}}
}
