package delete

import (
	"auth/pkg/model/json"
	resource_delete_service "auth/pkg/service/resource/delete"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db         *gorm.DB
	resourceId string
	fields     logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, resourceId string) *Api {
	return &Api{
		resourceId: resourceId,
		db:         db,
		fields:     logrus.Fields{"uuid": uuid, "url": url, "resource_id": resourceId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.resourceId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty resource_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty resource_id")
	}
	var resourceId int64
	if resourceId, err = strconv.ParseInt(a.resourceId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid resource_id")
	}
	httpStatus, resource, err := resource_delete_service.NewService(a.db, resourceId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"resource": resource}}
}
