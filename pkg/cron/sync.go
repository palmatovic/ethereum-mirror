package cron

import (
	database2 "ethereum-mirror/pkg/database"
	"ethereum-mirror/pkg/service/transaction"
	"ethereum-mirror/pkg/service/transaction_detail"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
)

type Env struct {
	playwright.Browser
	Database *gorm.DB
	Address  string
}

func (e *Env) SyncTransactions() (response interface{}, err error) {

	transactions, err := transaction.GetByLimitAndAddress(e.Browser, e.Address)
	if err != nil {
		return nil, err
	}

	if transactions != nil && len(transactions) > 0 {
		var savedTransactions []database2.Transaction
		if savedTransactions, err = transaction.SaveNew(e.Database, e.Address, transactions); err != nil {
			return nil, err
		}
		if savedTransactions != nil && len(savedTransactions) > 0 {
			_, err = transaction_detail.GetByTransaction(e.Browser, savedTransactions)
			if err != nil {
				return nil, err
			}

		}
	}

	return nil, nil
}
