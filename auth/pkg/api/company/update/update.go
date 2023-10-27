package update

import (
	company_db "auth/pkg/database/company"
	util_json "auth/pkg/model/json"
	company_update_service "auth/pkg/service/company/update"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Api struct {
	db     *gorm.DB
	body   []byte
	fields logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, body []byte) *Api {
	return &Api{
		body:   body,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "body": string(body)},
	}
}

func (a *Api) Update() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var company company_db.Company
	if err := json.Unmarshal(a.body, &company); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, companyDb, err := company_update_service.NewService(a.db, &company).Update()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"company": companyDb}}
}
