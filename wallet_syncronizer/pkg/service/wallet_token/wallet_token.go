package wallet_token

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	token_db "wallet-syncronizer/pkg/database/token"
	wallet_db "wallet-syncronizer/pkg/database/wallet"
	wallet_token_db "wallet-syncronizer/pkg/database/wallet_token"
	"wallet-syncronizer/pkg/service/alchemy/wallet_balance"
	token_find_or_create_service "wallet-syncronizer/pkg/service/token/find_or_create"
	string2 "wallet-syncronizer/pkg/util/string"
)

// FindOrCreateWalletTokens returns a list of all token balance by wallet_token
func FindOrCreateWalletTokens(walletDb wallet_db.Wallet, db *gorm.DB, alchemyApiKey string) (walletTokens []wallet_token_db.WalletToken, err error) {

	alchemyWalletBalances, err := wallet_balance.GetWalletBalances(string(walletDb.WalletId), alchemyApiKey)
	if err != nil {
		return nil, err
	}

	var (
		concurrentGoroutines = 10
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
		mutex                sync.Mutex
	)

	for i := range alchemyWalletBalances.Result.WalletBalances {
		if !string2.EmptyTokenBalance(alchemyWalletBalances.Result.WalletBalances[i].TokenBalance) {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(walletBalance wallet_balance.Balance) {
				defer func() {
					wg.Done()
					<-semaphore
				}()
				var logFields = logrus.Fields{"wallet_id": walletDb.WalletId, "token_id": walletBalance.TokenContractAddress}
				var tokenDb *token_db.Token
				tokenDb, err = token_find_or_create_service.NewService(db, walletBalance.TokenContractAddress, alchemyApiKey).FindOrCreateToken()
				if err != nil {
					logrus.WithFields(logFields).WithError(err).Error("cannot find or create token")
					return
				}
				//if tokenDb.RiskScam != 0 || tokenDb.WarningScam != 0 {
				//	//logrus.WithFields(logrus.Fields{"wallet_id": walletDb.WalletId, "token_contract_address": walletBalance.TokenContractAddress}).Warningf("scam token found. skipping process")
				//	return
				//}
				tokenAmount := string2.CalculateAmount(walletBalance.TokenBalance, tokenDb.Decimals)
				var walletToken wallet_token_db.WalletToken
				if errFirst := db.Where("WalletId = ? AND TokenId = ?", walletDb.WalletId, tokenDb.TokenId).First(&walletToken).Error; errFirst != nil {
					if errors.Is(errFirst, gorm.ErrRecordNotFound) {
						walletToken = wallet_token_db.WalletToken{
							WalletId:       walletDb.WalletId,
							TokenId:        tokenDb.TokenId,
							TokenAmount:    tokenAmount,
							TokenAmountHex: walletBalance.TokenBalance,
						}
						if errCreate := db.Create(&walletToken).Error; errCreate != nil {
							logrus.WithFields(logFields).WithError(errCreate).Error("cannot create wallet token")
						}
						mutex.Lock()
						walletTokens = append(walletTokens, walletToken)
						mutex.Unlock()
					} else {
						logrus.WithFields(logFields).WithError(errFirst).Error("cannot query wallet token")
						return
					}
				} else {
					if tokenAmount != walletToken.TokenAmount {
						walletToken.TokenAmount = tokenAmount
						if errUpdate := db.Updates(&walletToken).Error; errUpdate != nil {
							logrus.WithFields(logFields).WithError(errUpdate).Error("cannot update wallet token")
						}
						mutex.Lock()
						walletTokens = append(walletTokens, walletToken)
						mutex.Unlock()
					}
				}

			}(alchemyWalletBalances.Result.WalletBalances[i])
		}
	}
	wg.Wait()
	return walletTokens, nil
}
