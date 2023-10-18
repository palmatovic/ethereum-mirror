package get

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wallet-synchronizer/pkg/model/json"
	wallet_transaction_get_service "wallet-synchronizer/pkg/service/wallet_transaction/get"
)

type Api struct {
	db                  *gorm.DB
	walletTransactionId string
	fields              logrus.Fields
}

func NewApi(uuid string, url string, db *gorm.DB, walletTransactionId string) *Api {
	return &Api{
		walletTransactionId: walletTransactionId,
		db:                  db,
		fields:              logrus.Fields{"uuid": uuid, "url": url, "wallet_transaction_id": walletTransactionId},
	}
}

func (a *Api) Get() (status int, response interface{}) {
	logrus.WithFields(a.fields).Info("started")
	if len(a.walletTransactionId) == 0 {
		return fiber.StatusBadRequest, json.NewErrorResponse(fiber.StatusBadRequest, "empty token_id")
	}
	httpStatus, walletTransaction, err := wallet_transaction_get_service.NewService(a.db, a.walletTransactionId).Get()
	if err != nil {
		logrus.WithFields(a.fields).WithError(err).Errorf("terminated with failure")
		return httpStatus, json.NewErrorResponse(httpStatus, err.Error())
	}
	logrus.WithFields(a.fields).Info("terminated with success")
	return httpStatus, json.Response{Data: fiber.Map{"wallet_transaction": walletTransaction}}
}
