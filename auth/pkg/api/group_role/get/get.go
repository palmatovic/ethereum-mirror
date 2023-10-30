package get

import (
	"auth/pkg/model/json"
	group_role_get_service "auth/pkg/service/group_role/get"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db          *gorm.DB
	groupRoleId string
	fields      logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, groupRoleId string) *Api {
	return &Api{
		groupRoleId: groupRoleId,
		db:          db,
		fields:      logrus.Fields{"uuid": uuid, "url": url, "group_role_id": groupRoleId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.groupRoleId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty group_role_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty group_role_id")
	}
	var groupRoleId int64
	if groupRoleId, err = strconv.ParseInt(a.groupRoleId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid group_role_id")
	}
	httpStatus, groupRole, err := group_role_get_service.NewService(a.db, groupRoleId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"group_role": groupRole}}
}
