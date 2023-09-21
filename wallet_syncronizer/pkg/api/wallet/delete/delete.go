package delete

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	wallet_delete_service "wallet-syncronizer/pkg/service/wallet/delete"
	"wallet-syncronizer/pkg/util/json"
)

type Api struct {
	db       *gorm.DB
	walletId string
	fields   logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, walletId string) *Api {
	return &Api{
		walletId: walletId,
		db:       db,
		fields:   logrus.Fields{"uuid": uuid, "url": url, "wallet_id": walletId},
	}
}

func (a *Api) Delete() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.walletId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty wallet_id")
	}
	httpStatus, wallet, err := wallet_delete_service.NewService(a.db, a.walletId).Delete()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"wallet": wallet}}
}
