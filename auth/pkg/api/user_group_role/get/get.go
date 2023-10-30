package get

import (
	"auth/pkg/model/json"
	user_group_role_get_service "auth/pkg/service/user_group_role/get"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db              *gorm.DB
	userGroupRoleId string
	fields          logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, userGroupRoleId string) *Api {
	return &Api{
		userGroupRoleId: userGroupRoleId,
		db:              db,
		fields:          logrus.Fields{"uuid": uuid, "url": url, "user_group_role_id": userGroupRoleId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.userGroupRoleId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty user_group_role_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty user_group_role_id")
	}
	var userGroupRoleId int64
	if userGroupRoleId, err = strconv.ParseInt(a.userGroupRoleId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid user_group_role_id")
	}
	httpStatus, userGroupRole, err := user_group_role_get_service.NewService(a.db, userGroupRoleId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"user_group_role": userGroupRole}}
}
