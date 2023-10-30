package delete

import (
	"auth/pkg/model/json"
	role_delete_service "auth/pkg/service/role/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db     *gorm.DB
	roleId string
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, roleId string) *Api {
	return &Api{
		roleId: roleId,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "role_id": roleId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.roleId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty role_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty role_id")
	}
	var roleId int64
	if roleId, err = strconv.ParseInt(a.roleId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid role_id")
	}
	httpStatus, role, err := role_delete_service.NewService(a.db, roleId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"role": role}}
}
