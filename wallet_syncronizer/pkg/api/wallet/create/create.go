package create

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_db "wallet-syncronizer/pkg/database/wallet"
	wallet_create_service "wallet-syncronizer/pkg/service/wallet/create"
	util_json "wallet-syncronizer/pkg/util/json"
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

func (a *Api) Create() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var wallet wallet_db.Wallet
	if err := json.Unmarshal(a.body, &wallet); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, walletDb, err := wallet_create_service.NewService(a.db, &wallet).Create()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"wallet": walletDb}}
}
