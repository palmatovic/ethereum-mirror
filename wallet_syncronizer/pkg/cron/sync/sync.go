package sync

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	wallet_service "wallet-syncronizer/pkg/service/wallet"
	wallet_token_service "wallet-syncronizer/pkg/service/wallet_token"
)

type Env struct {
	Browser       playwright.Browser
	Database      *gorm.DB
	Wallets       []string
	AlchemyApiKey string
}

func (e *Env) Sync() {
	var (
		concurrentGoroutines = 10
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
	)

	for _, walletAddress := range e.Wallets {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(wAddress string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			wallet, err := wallet_service.FindOrCreateWallet(wAddress, e.Database)
			if err != nil {
				logrus.WithError(err).Error("cannot find or create wallet")
				return
			}

			walletTokens, err := wallet_token_service.FindOrCreateWalletTokens(wallet, e.Database, e.AlchemyApiKey)
			if err != nil {
				logrus.WithError(err).Error("cannot find or create wallet tokens")
				return
			}

			logrus.Infof("wallet tokens to update %d", len(walletTokens))

			//_, err = wallet_transaction_service.FindOrCreateWalletTransactions(walletTokensDb, e.Database, e.Browser)
			//if err != nil {
			//	logrus.WithError(err).Error("cannot find or create wallet transactions")
			//	return
			//}

		}(walletAddress)
	}
	wg.Wait()
	logrus.Info("sync terminated")
}
