package get

import (
	"auth/pkg/model/json"
	company_get_service "auth/pkg/service/company/get"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db        *gorm.DB
	companyId string
	fields    logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, companyId string) *Api {
	return &Api{
		companyId: companyId,
		db:        db,
		fields:    logrus.Fields{"uuid": uuid, "url": url, "company_id": companyId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var err error
	if len(a.companyId) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty company_id")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty company_id")
	}
	var companyId int64
	if companyId, err = strconv.ParseInt(a.companyId, 10, 64); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid company_id")
	}
	httpStatus, company, err := company_get_service.NewService(a.db, companyId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"company": company}}
}
