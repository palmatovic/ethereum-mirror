package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	syncronize "sync"
	"time"
	token_api "wallet-syncronizer/pkg/api/token"
	wallet_api "wallet-syncronizer/pkg/api/wallet"
	wallet_token_api "wallet-syncronizer/pkg/api/wallet_token"
	wallet_transaction_api "wallet-syncronizer/pkg/api/wallet_transaction"
	"wallet-syncronizer/pkg/cron/sync"
	"wallet-syncronizer/pkg/database/token"
	"wallet-syncronizer/pkg/database/wallet"
	"wallet-syncronizer/pkg/database/wallet_token"
	"wallet-syncronizer/pkg/database/wallet_transaction"
	token_url "wallet-syncronizer/pkg/util/url/token"
	wallet_url "wallet-syncronizer/pkg/util/url/wallet"
	wallet_token_url "wallet-syncronizer/pkg/util/url/wallet_token"
	wallet_transaction_url "wallet-syncronizer/pkg/util/url/wallet_transaction"
)

// Environment
// Defines the structure for holding environment variables
type Environment struct {
	//Wallets            []string `env:"WALLETS,required"` // 0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367
	PlaywrightHeadLess   bool   `env:"PLAYWRIGHT_HEADLESS,required"`
	AlchemyApiKey        string `env:"ALCHEMY_API_KEY,required"` // owUCVigVvnHA63o0C6mh3yrf3jxMkV7b
	FiberPort            int    `env:"FIBER_PORT,required"`
	BrowserPath          string `env:"BROWSER_PATH,required"`
	SrapeIntervalMinutes int    `env:"SCRAPE_INTERVAL_MINUTES,required"`
}

// WALLETS=0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367;PLAYWYRIGHT_HEADLESS=false;ALCHEMY_API_KEY=owUCVigVvnHA63o0C6mh3yrf3jxMkV7b

func main() {
	var (
		e   = Environment{}
		db  *gorm.DB
		err error
	)

	logrus.New()

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
	})
	logrus.SetReportCaller(true)

	// Parse environment variables into the 'e' struct
	if err = env.Parse(&e); err != nil {
		logrus.WithError(err).Fatalln("error during environment parsing")
	}

	// Open a connection to the SQLite database
	if db, err = gorm.Open(sqlite.Open("wallet.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		logrus.WithError(err).Fatalln("error during database connection")
	}

	// Perform automatic database schema migration
	err = db.AutoMigrate(&wallet.Wallet{}, &token.Token{}, &wallet_token.WalletToken{}, &wallet_transaction.WalletTransaction{})
	if err != nil {
		logrus.WithError(err).Fatalln("error during migration of database")
	}

	go cronJob(db, e.AlchemyApiKey, e.SrapeIntervalMinutes, e.BrowserPath, e.PlaywrightHeadLess)

	app := fiber.New()

	tokenApi := token_api.NewEnv(db)
	walletApi := wallet_api.NewEnv(db)
	walletTokenApi := wallet_token_api.NewEnv(db)
	walletTransactionApi := wallet_transaction_api.NewEnv(db)

	app.Get(token_url.Get, tokenApi.Get)
	app.Get(token_url.GetList, tokenApi.GetList)

	app.Get(wallet_url.Get, walletApi.Get)
	app.Get(wallet_url.GetList, walletApi.GetList)
	app.Post(wallet_url.Create, walletApi.Create)
	app.Put(wallet_url.Update, walletApi.Update)
	app.Delete(wallet_url.Delete, walletApi.Delete)

	app.Get(wallet_token_url.Get, walletTokenApi.Get)
	app.Get(wallet_token_url.GetList, walletTokenApi.GetList)

	app.Get(wallet_transaction_url.Get, walletTransactionApi.Get)
	app.Get(wallet_transaction_url.GetList, walletTransactionApi.GetList)

	if err = app.Listen(fmt.Sprintf(":%d", e.FiberPort)); err != nil {
		log.Fatalf("cannot starting fiber server: %v", err)
	}
}

func cronJob(db *gorm.DB, alchemyApiKey string, interval int, browserPath string, headless bool) {

	// Set up Playwright
	pw, err := playwright.Run()
	if err != nil {
		logrus.Fatalln("error during Playwright startup:", err)
	}

	defer func(pw *playwright.Playwright) {
		_ = pw.Stop()
	}(pw)

	// Install Playwright
	if err = playwright.Install(); err != nil {
		logrus.Fatalln("error during Playwright installation:", err)
	}

	// Launch Firefox browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(browserPath),
		Headless:       playwright.Bool(headless),
	})
	if err != nil {
		logrus.Fatalln("error during browser launch:", err)
	}

	// Create a mutex for synchronization
	var mutex syncronize.Mutex

	// Create an instance of the cron environment
	c := sync.Env{Browser: browser, Database: db, AlchemyApiKey: alchemyApiKey}

	// Create a new cron scheduler
	s := gocron.NewScheduler(time.Local)

	for {

		// Define the cron job using cron syntax
		_, err = s.Every(interval).Minute().Do(func() {
			// Lock the mutex before starting the task
			mutex.Lock()
			defer mutex.Unlock() // Unlock the mutex when the function finishes
			c.Sync()
		})

		if err != nil {
			logrus.WithError(err).Fatalln("cannot start cron")
		}

		// Start the cron scheduler (blocking call)
		s.StartBlocking()

		logrus.Warning("scheduler stopped, shutting down...restarting")
	}
}
