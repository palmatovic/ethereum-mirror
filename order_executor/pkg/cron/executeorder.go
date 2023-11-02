package sync

import (
	"errors"
	"fmt"
	db_t "order-executor/pkg/model/database/token"
	db_w "order-executor/pkg/model/database/wallet"
	db_wto "order-executor/pkg/model/database/wallet_token"
	db_wtr "order-executor/pkg/model/database/wallet_transaction"
	u "order-executor/pkg/util"
	"time"

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

func (e *Env) ExecuteOrders() (response interface{}, err error) {

	var personalWallet db_w.Wallet
	err = e.Database.Where("Type = ?", false).First(&personalWallet).Error
	if err != nil {
		return nil, err
	}

	// prendi tutte le transazioni recuperate non ancora processate

	var cryptoTransactions []db_wtr.WalletTransaction
	err = e.Database.Where("ProcessedByOrderExecutor = ?", false).Preload("Wallet").Find(&cryptoTransactions).Error

	if err != nil {
		return nil, err
	}

	if len(cryptoTransactions) > 0 {
		currentProcessDate := time.Now()

		for _, ct := range cryptoTransactions {

			// 	// verifica che la transazione non sia scaduta
			// qui va usata non la data di creazione ma la data della transazione.
			transactionExpired := isTransactionExpired(ct.AgeTimestamp, currentProcessDate, e.OrderTimeExpirationThreshold)
			if transactionExpired {
				// scaduta
				// marcala solo come processata
				// forse va aggiunto un campo di stato di processamento scaduta, aperta, chiusa, annullata per superamento soglia ecc
				err = setTransactionStatus(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_EXPIRED)
				if err != nil {
					// lancia errore e continua
				}
			}
			// prendo il wallet e poi il wallet token amount () usa TokenId di transactio.token per combaciare con wallet_token.

			if ct.TxType == u.WALLET_TRANSACTION_TYPE_BUY {
				// calcolo l'ammontare del wallet che sto monitorando
				totalWalletAmount, err := calculateTotalWalletAmount(e.Database, ct.WalletId, &ct.WalletTransactionId)
				totalTransactionAmount := ct.Price * ct.Amount
				percTransactionAmount := totalTransactionAmount * 100 / *totalWalletAmount

				// ora verifica le regole sulle percentuali.
				err = e.transactionPercentageCheck(&percTransactionAmount)

				// calcolo l'ammontare del mio wallet
				totalPersonalWalletAmount, err := calculateTotalWalletAmount(e.Database, personalWallet.WalletId, &ct.WalletTransactionId)
				totalPersonalTransactionAmount := *totalPersonalWalletAmount * *totalPersonalWalletAmount / 100
				fmt.Print(totalPersonalTransactionAmount)

				// totalPersonalTransactionAmount va diviso per il valore del token per sapere la quantità da comprare.
				// immaginiamo che sia uno

				// eseguo l'operazione e setto lo stato
				err = setTransactionStatus(e.Database, ct, false, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_OPENED)
				if err != nil {
					// lancia errore e continua
				}

				// prendo il token generico
				// da cui prenderò il valore corrente del token
				token := new(db_t.Token)
				err = e.Database.Where("TokenId = ?", ct.TokenId).First(token).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {

					} else {

					}
				}

				// prendo il token generico
				// da cui prenderò il valore dell'ammontare per il token.
				// se l'ammontare è zero alora
				walletToken := new(db_wto.WalletToken)
				err = e.Database.Where("TokenId = ? and WalletId = ?", ct.TokenId, ct.WalletId).First(walletToken).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {

					} else {

					}
				}

			} else if ct.TxType == u.WALLET_TRANSACTION_TYPE_SELL {

			} else {

			}
			// 	// prendi il saldo del wallet_token da cui parte la transazione,
			// 	// calcola la percentuale

			// il mmio saldo lo prendo mio wallet ()  c'è il type, se type = 0 allora è il mio. type 1 sono gli altri.
			// per i compra, posso non avere quel token nel mio wallet.
			// quindi devo comprare da eth verso il token di destinazione
			// devo capirese lui va da A --> B, io devo andare da ETH a B, quindi devo capire quanti ETH vale il suo cambio.  ho l'ammontare in dollari, allora calcolo da dollari ad eth. mi serve un'api
			//

			// se vedo una transazione in vendità, verifico se ho il token che lui sta vendendo, se non ce lo ho allora nulla,
			// se ce lo ho allora lo vendo

			// quindi lui compra X allora io compro da eth ad X.
			// quando lui vende, se ho quel token, lo vendo. se non ce lo ho passo
			// quando invece vendo io perchè mi basta, allora mi serve il valore in dollari corrente del token.

			// limitare il numero di transazioni attive e vendere quando hai un tot di guadagno, ma se lui vende prima del guadagno allora vendi
			// stessa cosa per la perdita.
			// 	// calcola il proprio saldo sullo stesso token.
			// se non ho fondi in quel token salto
			// 	// calcola la percentuale riportata al proprio saldo

			// 	// verifica che rispetti lee regole di threshold
		}
	}

	// genera il corrispettivo ordine

	// salva la transazione come processata

	return nil, nil
}

func isTransactionExpired(transactionDate time.Time, currentProcessDate time.Time, OrderTimeExpirationThreshold int) bool {

	transactionIsExpired := false
	expiredThresholdParameter := time.Duration(OrderTimeExpirationThreshold) * time.Minute
	expiredThresholdDate := currentProcessDate.Add(-expiredThresholdParameter)
	if transactionDate.Before(expiredThresholdDate) {
		transactionIsExpired = true
	}
	return transactionIsExpired
}

func setTransactionStatus(db *gorm.DB, transaction db_wtr.WalletTransaction, processed bool, status string) (err error) {

	transaction.ProcessedByOrderExecutor = processed
	if processed {
		transaction.ProcessedByOrderExecutorAt = time.Now()
	}
	transaction.ProcessedByOrderExecutorStatus = status
	err = db.Where("WalletTransactionId = ?", transaction.WalletTransactionId).Updates(&transaction).Error

	if err != nil {
		return err
	}
	return nil
}

func (e *Env) transactionPercentageCheck(percentage *float64) (err error) {

	maxPerc := float64(e.MaxPercOrderThreshold)
	minPerc := float64(e.MinPercOrderThreshold)

	if *percentage > maxPerc && e.SetMaxPercThreshold {
		percentage = &maxPerc
	} else if *percentage > maxPerc && !e.SetMaxPercThreshold {
		return errors.New("percentage too hight")
	} else if *percentage < minPerc && e.SetMinPercThreshold {
		percentage = &minPerc
	} else if *percentage < minPerc && !e.SetMinPercThreshold {
		return errors.New("percentage too low")
	}
	return nil
}

func calculateTotalWalletAmount(db *gorm.DB, walletId string, currentTransactionId *string) (result *float64, err error) {

	// per un wallett devo prendere tutti gli amount dei vari token e moltiplicarli per
	// il valore in dollari
	// e prendere tutte le transazioni successive alla data di aggiornamento dei wallet token
	// e sommare/sottrarre il loro amounth
	var totalAmount float64 = 0
	walletTokens := new([]db_wto.WalletToken)
	err = db.Where("WalletId = ?", walletId).Find(walletTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// in questo caso allora ritorno 0
			return &totalAmount, nil
		} else {

		}
		return nil, err
	}

	// cicla sui wallet
	// per ogni wallet prendi l'amount e le transazioni dopo la data di modifica
	if walletTokens != nil && len(*walletTokens) > 0 {
		for _, wto := range *walletTokens {
			// qui va moltiplicato il valore del token per il valore in dollari, per ora metto 1
			totalAmount = totalAmount + (wto.TokenAmount * 1)
			walletTransactions := new([]db_wtr.WalletTransaction)
			if currentTransactionId != nil {
				err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ? and TransactionId != ?", walletId, wto.TokenId, wto.UpdatedAt, currentTransactionId).Find(walletTransactions).Error
			} else {
				err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ?", walletId, wto.TokenId, wto.UpdatedAt).Find(walletTransactions).Error
			}

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {

				} else {

				}
				return nil, err
			}
			if walletTransactions != nil && len(*walletTransactions) > 0 {
				for _, wtr := range *walletTransactions {
					if wtr.TxType == u.WALLET_TRANSACTION_TYPE_BUY {
						totalAmount = totalAmount + (wtr.Price * wtr.Amount)
					} else if wtr.TxType == u.WALLET_TRANSACTION_TYPE_SELL {
						totalAmount = totalAmount - (wtr.Price * wtr.Amount)
					} else {
						totalAmount = totalAmount
					}
				}
			}
		}
	}
	result = &totalAmount
	return result, nil
}
