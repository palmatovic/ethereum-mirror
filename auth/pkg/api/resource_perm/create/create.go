package create

import (
	resource_perm_db "auth/pkg/database/resource_perm"
	util_json "auth/pkg/model/json"
	resource_perm_create_service "auth/pkg/service/resource_perm/create"
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
	var resourcePerm resource_perm_db.ResourcePerm
	if err := json.Unmarshal(a.body, &resourcePerm); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, resourcePermDb, err := resource_perm_create_service.NewService(a.db, &resourcePerm).Create()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"resource_perm": resourcePermDb}}
}
