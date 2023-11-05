package sync

import (
	"errors"
	"fmt"
	"math"
	db_t "order-executor/pkg/model/database/token"
	db_w "order-executor/pkg/model/database/wallet"
	db_wto "order-executor/pkg/model/database/wallet_token"
	db_wtr "order-executor/pkg/model/database/wallet_transaction"
	u "order-executor/pkg/util"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Env struct {
	MinPercOrderThreshold        int
	MaxPercOrderThreshold        int
	MinAbsOrderThreshold         int
	MaxAbsOrderThreshold         int
	SetMaxPercThreshold          bool
	SetMinPercThreshold          bool
	OrderTimeExpirationThreshold int
	StopEarningThreshold         int
	StopLossThreshold            int
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
	err = e.Database.Where("ProcessedByOrderExecutor = ?", false).Find(&cryptoTransactions).Error

	if err != nil {
		return nil, err
	}

	if len(cryptoTransactions) > 0 {
		currentProcessDate := time.Now()

		for _, ct := range cryptoTransactions {
			if ct.TxType == u.WALLET_TRANSACTION_TYPE_BUY {
				err = e.manageBuyTransaction(personalWallet, ct, currentProcessDate)
				if err != nil {
					// lancia errore e continua
					continue
				}
			} else if ct.TxType == u.WALLET_TRANSACTION_TYPE_SELL {
				err = e.manageSellTransaction(personalWallet, ct, currentProcessDate)
				if err != nil {
					// lancia errore e continua
					continue
				}
			} else {

			}

			// prendo il wallet e poi il wallet token amount () usa TokenId di transactio.token per combaciare con wallet_token.

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
	} else {
		logrus.Infof("no transctions to process")
	}

	// genera il corrispettivo ordine

	// salva la transazione come processata

	return nil, nil
}

// da finire
func (e *Env) manageSellTransaction(personalWallet db_w.Wallet, ct db_wtr.WalletTransaction, currentProcessDate time.Time) (err error) {
	// se vende, non faccio il conto di quanto ha venduto.
	// vendo tutto di quel token se ce lo ho
	if ct.ProcessedByOrderExecutorStatus == u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_OPENED {
		// verifica che non ci sia stata già una transazione di sell eseguita che abbia sovrascritto questa buy
		// per farlo basta hai ancora token associati al mio wallet
		var personalWalletTokenAbs *float64
		var personalWalletTokenAmount *float64
		personalWalletTokenAbs, personalWalletTokenAmount, err = calculateTotalWalletTokenAbs(e.Database, personalWallet.WalletId, ct.TokenId, nil)
		if err != nil {
			logrus.WithError(err).Errorf("failed to calculate personal wallet token abs")
			return err
		}
		if *personalWalletTokenAbs > 0 && *personalWalletTokenAmount > 0 {
			// ho ancora dei token da gestire e vendere
			// chiudi la transazione vendendo
			err = executeOperation(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_CLOSED)
			if err != nil {
				logrus.WithError(err).Errorf("failed to execute operation on chain")
				return err
			}
		} else {
			// non ci sono più token da vendere quindi
			// devo chiudere la transazione impostando il suo stato
			err = setTransactionStatus(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_ALREADY_CLOSED)
			if err != nil {
				logrus.WithError(err).Errorf("failed to set transaction status")
				return err
			}
		}
	}
	return nil
}

func (e *Env) manageBuyTransaction(personalWallet db_w.Wallet, ct db_wtr.WalletTransaction, currentProcessDate time.Time) (err error) {
	// se è registrata devo aprirla comprando
	// se è aperta devo capire se è ora di chiuderla e vendere
	if ct.ProcessedByOrderExecutorStatus == u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_REGISTERED {
		// verifica che la transazione non sia scaduta
		transactionExpired := isTransactionExpired(ct.AgeTimestamp, currentProcessDate, e.OrderTimeExpirationThreshold)
		if transactionExpired {
			// scaduta
			// marcala solo come processata
			// forse va aggiunto un campo di stato di processamento scaduta, aperta, chiusa, annullata per superamento soglia ecc
			err = setTransactionStatus(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_EXPIRED)
			if err != nil {
				// lancia errore e continua
				logrus.WithError(err).Errorf("failed to set transaction status")
				return err
			}
		}

		// se non è scaduta
		// calcolo l'ammontare del wallet che sto monitorando
		totalWalletAbs, err := calculateTotalWalletAbs(e.Database, ct.WalletId, &ct.WalletTransactionId)
		if err != nil {
			// lancia errore e continua
			logrus.WithError(err).Errorf("failed to set calculate monitored wallet amounth")
			return err
		}
		totalTransactionAbs := ct.Price * ct.Amount
		totalTransactionPerc := totalTransactionAbs * 100 / *totalWalletAbs

		// ora verifica le regole sulle percentuali.
		// ovvero se non vado sopra la massima percentuale o sotto la minima
		err = e.transactionPercentageCheck(&totalTransactionPerc)
		if err != nil {
			// lancia errore e continua
			logrus.WithError(err).Errorf("transaction percentage check failed")
			return err
		}

		// calcolo la cifra in dollari/euro del mio wallet
		totalPersonalWalletAbs, err := calculateTotalWalletAbs(e.Database, personalWallet.WalletId, &ct.WalletTransactionId)
		if err != nil {
			// lancia errore e continua
			logrus.WithError(err).Errorf("failed to set calculate personal wallet amounth")
			return err
		}
		totalPersonalTransactionAbs := totalTransactionPerc * *totalPersonalWalletAbs / 100

		// ora verifico se la cifra che emerge è superiore o inferiore al massimo /minimo in assoluto
		err = e.transactionAbsCheck(&totalPersonalTransactionAbs)
		if err != nil {
			// lancia errore e continua
			logrus.WithError(err).Errorf("transaction abs check failed")
			return err
		}

		// totalPersonalTransactionAmount va diviso per il valore del token per sapere la quantità da comprare.
		// immaginiamo che sia uno
		var tokenPrice *float64
		tokenPrice, err = getCurrentTokenAbs(e.Database, ct.TokenId)
		if err != nil {
			return err
		}
		numberOfTokensToBuy := totalPersonalTransactionAbs / *tokenPrice
		fmt.Print(numberOfTokensToBuy)

		// apro transazione ed eseguo prima l'inserimento su db e poi la chiamata alle api
		// se chiamata alle api va in errore effettuo il rollback
		err = executeOperation(e.Database, ct, false, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_OPENED)
		if err != nil {
			logrus.WithError(err).Errorf("failed to execute operation on chain")
			return err
		}
	} else if ct.ProcessedByOrderExecutorStatus == u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_OPENED {
		// verifica che non ci sia stata già una transazione di sell eseguita che abbia sovrascritto questa buy
		// per farlo basta hai ancora token associati al mio wallet
		var personalWalletTokenAbs *float64
		var personalWalletTokenAmount *float64
		personalWalletTokenAbs, personalWalletTokenAmount, err = calculateTotalWalletTokenAbs(e.Database, personalWallet.WalletId, ct.TokenId, nil)
		if err != nil {
			logrus.WithError(err).Errorf("failed to calculate personal wallet token abs")
			return err
		}
		if *personalWalletTokenAbs > 0 {
			// ho ancora dei token da gestire e vendere
			// prendi valore attuale del token rispetto all'apertura e verifica se va chiusa in base
			// a percentuale di perdita o guadagno impostati.
			var startTransactionAbs float64
			startTransactionAbs = ct.Price * ct.Amount
			var endTransactionAbs float64
			var tokenPrice *float64
			tokenPrice, err = getCurrentTokenAbs(e.Database, ct.TokenId)
			if err != nil {
				return err
			}
			endTransactionAbs = *tokenPrice
			delta := math.Abs(endTransactionAbs - startTransactionAbs)
			negativeSign := math.Signbit(endTransactionAbs - startTransactionAbs)
			perc := delta * 100 / startTransactionAbs
			if negativeSign && perc >= float64(e.StopLossThreshold) {
				// vendi e smetti di perdere
				fmt.Sprint(personalWalletTokenAmount)
				// apri transazione e vendilo
				err = executeOperation(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_CLOSED)
				if err != nil {
					logrus.WithError(err).Errorf("failed to execute operation on chain")
					return err
				}
			} else if !negativeSign && perc >= float64(e.StopEarningThreshold) {
				// vendi e smetti di guadagnare
				err = executeOperation(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_CLOSED)
				if err != nil {
					logrus.WithError(err).Errorf("failed to execute operation on chain")
					return err
				}
			} else {
				// continua a mantenere aperta l'operazione e non fare nulla
			}
		} else {
			// non ci sono più token da vendere quindi
			// devo chiudere la transazione impostando il suo stato
			err = setTransactionStatus(e.Database, ct, true, u.WALLET_TRANSACTION_ORDER_EXECUTOR_STATUS_ALREADY_CLOSED)
			if err != nil {
				logrus.WithError(err).Errorf("failed to set transaction status")
				return err
			}
		}
	}
	return nil
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
	updateDate := time.Now()
	transaction.ProcessedByOrderExecutor = processed
	if processed {
		transaction.ProcessedByOrderExecutorAt = updateDate
	}
	transaction.ProcessedByOrderExecutorStatus = status
	transaction.UpdatedAt = updateDate
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

func (e *Env) transactionAbsCheck(abs *float64) (err error) {

	maxAbs := float64(e.MaxPercOrderThreshold)
	minAbs := float64(e.MinPercOrderThreshold)

	if *abs > maxAbs && e.SetMaxPercThreshold {
		abs = &maxAbs
	} else if *abs > maxAbs && !e.SetMaxPercThreshold {
		return errors.New("abs too hight")
	} else if *abs < minAbs && e.SetMinPercThreshold {
		abs = &minAbs
	} else if *abs < minAbs && !e.SetMinPercThreshold {
		return errors.New("abs too low")
	}
	return nil
}

func calculateTotalWalletAbs(db *gorm.DB, walletId string, currentTransactionId *string) (result *float64, err error) {

	// per un wallett devo prendere tutti gli amount dei vari token e moltiplicarli per
	// il valore in dollari
	// e prendere tutte le transazioni successive alla data di aggiornamento dei wallet token
	// e sommare/sottrarre il loro amounth
	var totalAbs float64 = 0
	walletTokens := new([]db_wto.WalletToken)
	err = db.Where("WalletId = ?", walletId).Find(walletTokens).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// in questo caso allora ritorno 0
			return &totalAbs, nil
		} else {
			return nil, err
		}
	}

	// cicla sui wallet
	// per ogni wallet prendi l'amount e le transazioni dopo la data di modifica
	if walletTokens != nil && len(*walletTokens) > 0 {
		for _, wto := range *walletTokens {
			// qui va moltiplicato il valore del token per il valore in dollari, per ora metto 1
			tokenPrice, err := getCurrentTokenAbs(db, wto.TokenId)
			if err != nil {
				return nil, err
			}
			totalAbs = totalAbs + (wto.TokenAmount * *tokenPrice)
			walletTransactions := new([]db_wtr.WalletTransaction)
			if currentTransactionId != nil {
				err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ? and TransactionId != ?", walletId, wto.TokenId, wto.UpdatedAt, currentTransactionId).Find(walletTransactions).Error
			} else {
				err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ?", walletId, wto.TokenId, wto.UpdatedAt).Find(walletTransactions).Error
			}

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// non c'è problema se non ci sono transazioni dopo l'ultimo
					// aggiornamento del wallet
					result = &totalAbs
					return result, nil
				} else {
					return nil, err
				}
			}
			if walletTransactions != nil && len(*walletTransactions) > 0 {
				for _, wtr := range *walletTransactions {
					if wtr.TxType == u.WALLET_TRANSACTION_TYPE_BUY {
						totalAbs = totalAbs + (wtr.Price * wtr.Amount)
					} else if wtr.TxType == u.WALLET_TRANSACTION_TYPE_SELL {
						totalAbs = totalAbs - (wtr.Price * wtr.Amount)
					} else {
						totalAbs = totalAbs
					}
				}
			}
		}
	}
	result = &totalAbs
	return result, nil
}

func calculateTotalWalletTokenAbs(db *gorm.DB, walletId string, tokenId string, currentTransactionId *string) (*float64, *float64, error) {

	// per un wallett devo prendere gli amount di un solo token e moltiplicarlo per
	// il valore in dollari
	// e prendere tutte le transazioni successive alla data di aggiornamento dei wallet token
	// e sommare/sottrarre il loro amounth
	var err error
	var totalAbs float64 = 0
	var totalAmount float64 = 0
	walletToken := new(db_wto.WalletToken)
	err = db.Where("WalletId = ? and TokenId = ?", walletId, tokenId).First(walletToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// in questo caso allora ritorno 0
			return &totalAbs, &totalAmount, nil
		} else {
			return nil, nil, err
		}
	}

	// cicla sui wallet
	// per ogni wallet prendi l'amount e le transazioni dopo la data di modifica
	if walletToken != nil {
		// qui va moltiplicato il valore del token per il valore in dollari, per ora metto 1
		var tokenPrice *float64
		tokenPrice, err = getCurrentTokenAbs(db, tokenId)
		if err != nil {
			return nil, nil, err
		}
		totalAbs = totalAbs + (walletToken.TokenAmount * *tokenPrice)
		walletTransactions := new([]db_wtr.WalletTransaction)
		if currentTransactionId != nil {
			err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ? and TransactionId != ?", walletId, walletToken.TokenId, walletToken.UpdatedAt, currentTransactionId).Find(walletTransactions).Error
		} else {
			err = db.Where("WalletId = ? and TokenId = ? and AgeTimestamp > ?", walletId, walletToken.TokenId, walletToken.UpdatedAt).Find(walletTransactions).Error
		}

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &totalAbs, &totalAmount, nil
			} else {
				return nil, nil, err
			}
		}
		if walletTransactions != nil && len(*walletTransactions) > 0 {
			for _, wtr := range *walletTransactions {
				if wtr.TxType == u.WALLET_TRANSACTION_TYPE_BUY {
					totalAbs = totalAbs + (wtr.Price * wtr.Amount)
					totalAmount = totalAmount + wtr.Amount
				} else if wtr.TxType == u.WALLET_TRANSACTION_TYPE_SELL {
					totalAbs = totalAbs - (wtr.Price * wtr.Amount)
					totalAmount = totalAmount - wtr.Amount
				} else {
					totalAbs = totalAbs
					totalAmount = totalAmount
				}
			}
		}
	}
	return &totalAbs, &totalAmount, nil
}

func getCurrentTokenAbs(db *gorm.DB, tokenId string) (result *float64, err error) {

	tokenPrice := new(db_t.TokenPrice)
	err = db.Where("TokenId = ?", tokenId).Order("PriceDate DESC").First(tokenPrice).Error
	if err != nil {
		return nil, err
	}
	result = &((*tokenPrice).Price)
	return result, nil
}

// sicuramente devo passare anche il token l'amount dell'operazione
func executeOperation(db *gorm.DB, ct db_wtr.WalletTransaction, executorProcessBool bool, executorProcessStatus string) (err error) {

	tx := db.Begin()

	if tx.Error != nil {
		logrus.WithError(err).Errorf("failed to open db transaction")
		return err
	}

	// eseguo l'operazione e setto lo stato
	err = setTransactionStatus(db, ct, executorProcessBool, executorProcessStatus)
	if err != nil {
		tx.Rollback()
		logrus.WithError(err).Errorf("failed to set transaction status")
		return err
	}

	// EFFETTUA OPERAZIONE API
	// ...

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logrus.WithError(err).Errorf("failed to commit transaction")
		return err
	}
	return nil
}
