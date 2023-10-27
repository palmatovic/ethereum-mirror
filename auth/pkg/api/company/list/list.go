package get

import (
	"auth/pkg/model/json"
	company_list_service "auth/pkg/service/company/list"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type Api struct {
	db         *gorm.DB
	fields     logrus.Fields
	pageSize   string
	pageNumber string
}

func NewApi(
	uuid string,
	url string,
	db *gorm.DB,
	pageSize string,
	pageNumber string,
) *Api {
	return &Api{
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "page_size": pageSize, "page_number": pageNumber},
	}
}

func (a *Api) List() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var pageSize int
	var err error
	if len(a.pageSize) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty page_size")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty page_size")
	}

	if pageSize, err = strconv.Atoi(a.pageSize); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid page_size")
	}
	var pageNumber int
	if len(a.pageNumber) == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("empty page_number")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty page_number")
	}

	if pageNumber, err = strconv.Atoi(a.pageNumber); err != nil {
		logrus.WithFields(a.fields).WithError(err).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "invalid page_number")
	}
	if pageNumber == 0 {
		logrus.WithFields(a.fields).WithError(errors.New("page_number cannot be zero")).Error("terminated with failure")
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "page_number cannot be zero")
	}
	httpStatus, companys, err := company_list_service.NewService(a.db, pageSize, pageNumber).List()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"companys": companys}}
}
