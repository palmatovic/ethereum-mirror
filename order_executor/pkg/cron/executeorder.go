package sync

import (
	database "order-executor/pkg/model/database"

	"gorm.io/gorm"
)

type Env struct {
	MinPercOrderThreshold        int
	MaxPercOrderThreshold        int
	SetMaxPercThreshold          bool
	SetMinPercThreshold          bool
	OrderTimeExpirationThreshold int
	MaxPriceRangePerc            float32
	Database                     *gorm.DB
}

func (e *Env) ExecuteOrdres() (response interface{}, err error) {

	// prendi tutte le transazioni recuperate non ancora processate

	var cryptoTransactions []database.Transaction
	err = e.Database.Where("ProcessedByOrderExecutor = ?", false).Preload("FollowedWallet").Find(&cryptoTransactions).Error

	if err != nil {
		return nil, err
	}

	if len(cryptoTransactions) > 0 {
		// for _, ct := range cryptoTransactions {

		// 	// verifica che la transazione non sia scaduta
		// 	ct.
		// 	// prendi il saldo del wallet_token da cui parte la transazione,
		// 	// calcola la percentuale

		// 	// calcola il proprio saldo
		// 	// calcola la percentuale riportata al proprio saldo

		// 	// verifica che rispetti le
		// }
	}

	// genera il corrispettivo ordine

	// salva la transazione come processata

	return nil, nil
}
