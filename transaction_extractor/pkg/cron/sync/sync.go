package sync

import (
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	database2 "transaction-extractor/pkg/database/transaction"
	transaction2 "transaction-extractor/pkg/model/transaction"
	"transaction-extractor/pkg/service/transaction"
	"transaction-extractor/pkg/service/transaction_detail"
)

type Env struct {
	playwright.Browser
	Database  *gorm.DB
	Addresses []string
}

func (e *Env) SyncTransactions() (response interface{}, err error) {
	for _, address := range e.Addresses {
		var transactions []transaction2.Transaction
		transactions, err = transaction.GetByAddress(e.Browser, address)
		if err != nil {
			return nil, err
		}

		if transactions != nil && len(transactions) > 0 {
			var savedTransactions []database2.Transaction
			if savedTransactions, err = transaction.SaveNew(e.Database, transactions); err != nil {
				return nil, err
			}
			if savedTransactions != nil && len(savedTransactions) > 0 {
				_, err = transaction_detail.GetByTransaction(e.Browser, savedTransactions)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}
