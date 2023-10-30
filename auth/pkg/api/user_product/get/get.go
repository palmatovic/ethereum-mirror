package get

import (
	"auth/pkg/model/json"
	user_product_get_service "auth/pkg/service/user_product/get"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db            *gorm.DB
	userProductId string
	fields        logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, userProductId string) *Api {
	return &Api{
		userProductId: userProductId,
		db:            db,
		fields:        logrus.Fields{"uuid": uuid, "url": url, "user_product_id": userProductId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.userProductId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty user_product_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty user_product_id")
	}
	var userProductId int64
	if userProductId, err = strconv.ParseInt(a.userProductId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid user_product_id")
	}
	httpStatus, userProduct, err := user_product_get_service.NewService(a.db, userProductId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"user_product": userProduct}}
}
