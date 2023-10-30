package delete

import (
	"auth/pkg/model/json"
	user_delete_service "auth/pkg/service/user/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db     *gorm.DB
	userId string
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, userId string) *Api {
	return &Api{
		userId: userId,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "user_id": userId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.userId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty user_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty user_id")
	}
	var userId int64
	if userId, err = strconv.ParseInt(a.userId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid user_id")
	}
	httpStatus, user, err := user_delete_service.NewService(a.db, userId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"user": user}}
}
