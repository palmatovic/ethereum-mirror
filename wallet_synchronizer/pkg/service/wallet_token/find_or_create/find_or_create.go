package find_or_create

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	token_db "wallet-synchronizer/pkg/database/token"
	wallet_db "wallet-synchronizer/pkg/database/wallet"
	wallet_token_db "wallet-synchronizer/pkg/database/wallet_token"
	"wallet-synchronizer/pkg/service/alchemy/wallet_balance"
	token_find_or_create_service "wallet-synchronizer/pkg/service/token/find_or_create"
	wallet_token_create_service "wallet-synchronizer/pkg/service/wallet_token/create"
	wallet_token_get_service "wallet-synchronizer/pkg/service/wallet_token/get"
	wallet_token_update_service "wallet-synchronizer/pkg/service/wallet_token/update"

	util_string "wallet-synchronizer/pkg/util/string"
)

type Service struct {
	alchemyApiKey string
	walletDb      wallet_db.Wallet
	db            *gorm.DB
}

func NewService(db *gorm.DB, alchemyApiKey string, walletDb wallet_db.Wallet) *Service {
	return &Service{alchemyApiKey: alchemyApiKey, walletDb: walletDb, db: db}
}

// FindOrCreateWalletTokens returns a list of all token balance by wallet_token
func (s *Service) FindOrCreateWalletTokens() (walletTokens []wallet_token_db.WalletToken, err error) {

	alchemyWalletBalances, err := wallet_balance.GetWalletBalances(s.walletDb.WalletId, s.alchemyApiKey)
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
		if !util_string.EmptyTokenBalance(alchemyWalletBalances.Result.WalletBalances[i].TokenBalance) {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(walletBalance wallet_balance.Balance) {
				defer func() {
					wg.Done()
					<-semaphore
				}()
				var logFields = logrus.Fields{"wallet_id": s.walletDb.WalletId, "token_id": walletBalance.TokenContractAddress}
				var tokenDb *token_db.Token
				tokenDb, err = token_find_or_create_service.NewService(s.db, walletBalance.TokenContractAddress, s.alchemyApiKey).FindOrCreateToken()
				if err != nil {
					logrus.WithFields(logFields).WithError(err).Error("cannot find or create token")
					return
				}
				//if tokenDb.RiskScam != 0 || tokenDb.WarningScam != 0 {
				//	//logrus.WithFields(logrus.Fields{"wallet_id": walletDb.WalletId, "token_contract_address": walletBalance.TokenContractAddress}).Warningf("scam token found. skipping process")
				//	return
				//}
				tokenAmount := util_string.CalculateAmount(walletBalance.TokenBalance, tokenDb.Decimals)
				var walletToken *wallet_token_db.WalletToken

				if _, walletToken, err = wallet_token_get_service.NewService(s.db, s.walletDb.WalletId, tokenDb.TokenId).Get(); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						walletToken = &wallet_token_db.WalletToken{
							WalletId:       s.walletDb.WalletId,
							TokenId:        tokenDb.TokenId,
							TokenAmount:    tokenAmount,
							TokenAmountHex: walletBalance.TokenBalance,
						}
						_, _, err = wallet_token_create_service.NewService(s.db, walletToken).Create()
						if err != nil {
							logrus.WithFields(logFields).WithError(err).Error("cannot update wallet token")
						}
						mutex.Lock()
						walletTokens = append(walletTokens, *walletToken)
						mutex.Unlock()
					} else {
						logrus.WithFields(logFields).WithError(err).Error("cannot query wallet token")
						return
					}
				} else {
					if tokenAmount != walletToken.TokenAmount {
						walletToken.TokenAmount = tokenAmount
						_, _, err = wallet_token_update_service.NewService(s.db, walletToken).Update()
						if err != nil {
							logrus.WithFields(logFields).WithError(err).Error("cannot update wallet token")
						}
						mutex.Lock()
						walletTokens = append(walletTokens, *walletToken)
						mutex.Unlock()
					}
				}
			}(alchemyWalletBalances.Result.WalletBalances[i])
		}
	}
	wg.Wait()
	return walletTokens, nil
}
