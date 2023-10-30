package delete

import (
	"auth/pkg/model/json"
	product_delete_service "auth/pkg/service/product/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db        *gorm.DB
	productId string
	fields    logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, productId string) *Api {
	return &Api{
		productId: productId,
		db:        db,
		fields:    logrus.Fields{"uuid": uuid, "url": url, "product_id": productId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.productId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty product_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty product_id")
	}
	var productId int64
	if productId, err = strconv.ParseInt(a.productId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid product_id")
	}
	httpStatus, product, err := product_delete_service.NewService(a.db, productId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"product": product}}
}
