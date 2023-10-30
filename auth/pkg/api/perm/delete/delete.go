package delete

import (
	"auth/pkg/model/json"
	perm_delete_service "auth/pkg/service/perm/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db     *gorm.DB
	permId string
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, permId string) *Api {
	return &Api{
		permId: permId,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "perm_id": permId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.permId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty perm_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty perm_id")
	}
	var permId int64
	if permId, err = strconv.ParseInt(a.permId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid perm_id")
	}
	httpStatus, perm, err := perm_delete_service.NewService(a.db, permId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"perm": perm}}
}
