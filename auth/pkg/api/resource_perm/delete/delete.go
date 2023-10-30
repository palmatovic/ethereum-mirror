package delete

import (
	"auth/pkg/model/json"
	resource_perm_delete_service "auth/pkg/service/resource_perm/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db             *gorm.DB
	resourcePermId string
	fields         logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, resourcePermId string) *Api {
	return &Api{
		resourcePermId: resourcePermId,
		db:             db,
		fields:         logrus.Fields{"uuid": uuid, "url": url, "resource_perm_id": resourcePermId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.resourcePermId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty resource_perm_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty resource_perm_id")
	}
	var resourcePermId int64
	if resourcePermId, err = strconv.ParseInt(a.resourcePermId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid resource_perm_id")
	}
	httpStatus, resourcePerm, err := resource_perm_delete_service.NewService(a.db, resourcePermId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"resource_perm": resourcePerm}}
}
