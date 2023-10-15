package update

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
	util_json "wallet-synchronizer/pkg/model/json"
	wallet_update_service "wallet-synchronizer/pkg/service/wallet/update"
)

// Deprecated
type Api struct {
	db     *gorm.DB
	body   []byte
	fields logrus.Fields
}

// Deprecated
func NewApi(uuid string, url string, db *gorm.DB, body []byte) *Api {
	return &Api{
		body:   body,
		db:     db,
		fields: logrus.Fields{"uuid": uuid, "url": url, "body": string(body)},
	}
}

// Deprecated
func (a *Api) Update() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	var wallet wallet_db.Wallet
	if err := json.Unmarshal(a.body, &wallet); err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return fiber.StatusInternalServerError, util_json.NewErrorResponse(fiber.StatusInternalServerError, err.Error())
	}
	httpStatus, walletDb, err := wallet_update_service.NewService(a.db, &wallet).Update()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, util_json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, util_json.Response{Data: fiber.Map{"wallet": walletDb}}
}
