package get

import (
	"auth/pkg/model/json"
	group_get_service "auth/pkg/service/group/get"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db      *gorm.DB
	groupId string
	fields  logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, groupId string) *Api {
	return &Api{
		groupId: groupId,
		db:      db,
		fields:  logrus.Fields{"uuid": uuid, "url": url, "group_id": groupId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.groupId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty group_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty group_id")
	}
	var groupId int64
	if groupId, err = strconv.ParseInt(a.groupId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid group_id")
	}
	httpStatus, group, err := group_get_service.NewService(a.db, groupId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"group": group}}
}
