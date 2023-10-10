package sync

import (
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
	"wallet-synchronizer/pkg/database/wallet"
	wallet_find_or_create_service "wallet-synchronizer/pkg/service/wallet/find_or_create"
	wallet_token_service "wallet-synchronizer/pkg/service/wallet_token/find_or_create"
	wallet_transaction_service "wallet-synchronizer/pkg/service/wallet_transaction"
)

type Sync struct {
	Browser       playwright.Browser
	Database      *gorm.DB
	AlchemyApiKey string
}

func NewSync(browser playwright.Browser, db *gorm.DB, alchemyApiKey string) *Sync {
	return &Sync{
		Browser:       browser,
		Database:      db,
		AlchemyApiKey: alchemyApiKey,
	}
}

func (e *Sync) Sync() {
	logrus.Infof("sync started")

	var wallets []wallet.Wallet

	if err := e.Database.Find(&wallets).Error; err != nil {
		logrus.WithError(err).Errorf("failed to find wallets")
		return
	}
	var (
		concurrentGoroutines = 10
		semaphore            = make(chan struct{}, concurrentGoroutines)
		wg                   sync.WaitGroup
		startTime            = time.Now()
	)

	for _, w := range wallets {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(wAddress string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			wall, err := wallet_find_or_create_service.NewService(e.Database, wAddress).FindOrCreateWallet()
			if err != nil {
				logrus.WithError(err).Error("cannot find or create wallet")
				return
			}

			walletTokens, err := wallet_token_service.NewService(e.Database, e.AlchemyApiKey, *wall).FindOrCreateWalletTokens()
			if err != nil {
				logrus.WithError(err).Error("cannot find or create wallet tokens")
				return
			}

			err = wallet_transaction_service.FindOrCreateWalletTransactions(e.Database, walletTokens, e.Browser)
			if err != nil {
				logrus.WithError(err).Error("cannot find or create wallet transactions")
				return
			}

		}(w.WalletId)
	}
	wg.Wait()
	logrus.Infof("sync terminated in %s", time.Now().Sub(startTime).String())
}
