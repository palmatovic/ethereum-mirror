package create

import (
	resource_db "auth/pkg/database/resource"
	util_json "auth/pkg/model/json"
	resource_create_service "auth/pkg/service/resource/create"
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
		fields: logrus.Fields{"uuid": uuid, "url": url, "body": util_json.Stringify(body)},
	}
}

func (a *Api) Create() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var resource resource_db.Resource
	if err := json.Unmarshal(a.body, &resource); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, resourceDb, err := resource_create_service.NewService(a.db, &resource).Create()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"resource": resourceDb}}
}
