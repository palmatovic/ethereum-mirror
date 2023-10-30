package create

import (
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
	util_json "auth/pkg/model/json"
	group_role_resource_perm_create_service "auth/pkg/service/group_role_resource_perm/create"
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
	var groupRoleResourcePerm group_role_resource_perm_db.GroupRoleResourcePerm
	if err := json.Unmarshal(a.body, &groupRoleResourcePerm); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, groupRoleResourcePermDb, err := group_role_resource_perm_create_service.NewService(a.db, &groupRoleResourcePerm).Create()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"group_role_resource_perm": groupRoleResourcePermDb}}
}
