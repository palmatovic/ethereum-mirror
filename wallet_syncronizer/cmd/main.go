package main

import (
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	token_api "wallet-syncronizer/pkg/api/token"
	wallet_api "wallet-syncronizer/pkg/api/wallet"
	wallet_token_api "wallet-syncronizer/pkg/api/wallet_token"
	wallet_transaction_api "wallet-syncronizer/pkg/api/wallet_transaction"
	syncronizer "wallet-syncronizer/pkg/cron/sync"
	"wallet-syncronizer/pkg/database/token"
	"wallet-syncronizer/pkg/database/wallet"
	"wallet-syncronizer/pkg/database/wallet_token"
	"wallet-syncronizer/pkg/database/wallet_transaction"
	token_url "wallet-syncronizer/pkg/util/url/token"
	wallet_url "wallet-syncronizer/pkg/util/url/wallet"
	wallet_token_url "wallet-syncronizer/pkg/util/url/wallet_token"
	wallet_transaction_url "wallet-syncronizer/pkg/util/url/wallet_transaction"
)

type Environment struct {
	PlaywrightHeadLess    bool   `env:"PLAYWRIGHT_HEADLESS,required"`
	AlchemyApiKey         string `env:"ALCHEMY_API_KEY,required"`
	FiberPort             int    `env:"FIBER_PORT,required"`
	BrowserPath           string `env:"BROWSER_PATH,required"`
	ScrapeIntervalMinutes int    `env:"SCRAPE_INTERVAL_MINUTES,required"`
}

func main() {
	initializeLogger()
	e := loadEnvironment()
	db := initializeDatabase()
	initializeDatabaseSchema(db)
	go startCronJob(db, e.BrowserPath, e.PlaywrightHeadLess, e.AlchemyApiKey, e.ScrapeIntervalMinutes)
	app := initializeFiberApp(db)
	startFiberServer(app, e.FiberPort)
}

func loadEnvironment() Environment {
	var e Environment
	if err := env.Parse(&e); err != nil {
		logrus.WithError(err).Fatalln("error during environment parsing")
	}
	return e
}

func initializeLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
	})
	logrus.SetReportCaller(true)
}

func initializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("wallet.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		logrus.WithError(err).Fatalln("error during database connection")
	}

	return db
}

func initializeDatabaseSchema(db *gorm.DB) {
	err := db.AutoMigrate(&wallet.Wallet{}, &token.Token{}, &wallet_token.WalletToken{}, &wallet_transaction.WalletTransaction{})
	if err != nil {
		logrus.WithError(err).Fatalln("error during migration of database")
	}
	// TODO:
	// to be removed
	if err = db.Where("WalletId = ?", "0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367").First(&wallet.Wallet{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Create(&wallet.Wallet{WalletId: "0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367"}).Error; err != nil {
				logrus.WithError(err).Fatalln("error during initial wallet setup")
			}
		}
	}
}

func initializeFiberApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	// generate random request id for each call
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
		ContextKey: "uuid",
	}))

	app.Server().WriteTimeout = 300 * time.Second
	app.Server().ReadTimeout = 300 * time.Second
	app.Server().ReadBufferSize = 100 * 1024 * 1024
	app.Server().MaxRequestBodySize = 100 * 1024 * 1024

	tokenApi := token_api.NewEnv(db)
	walletApi := wallet_api.NewEnv(db)
	walletTokenApi := wallet_token_api.NewEnv(db)
	walletTransactionApi := wallet_transaction_api.NewEnv(db)

	app.Get(token_url.Get, tokenApi.Get)
	app.Get(token_url.List, tokenApi.List)
	app.Get(wallet_url.Get, walletApi.Get)
	app.Get(wallet_url.List, walletApi.List)
	app.Post(wallet_url.Create, walletApi.Create)
	app.Put(wallet_url.Update, walletApi.Update)
	app.Delete(wallet_url.Delete, walletApi.Delete)
	app.Get(wallet_token_url.Get, walletTokenApi.Get)
	app.Get(wallet_token_url.List, walletTokenApi.List)
	app.Get(wallet_transaction_url.Get, walletTransactionApi.Get)
	app.Get(wallet_transaction_url.GetList, walletTransactionApi.List)

	return app
}

func startFiberServer(app *fiber.App, port int) {
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("cannot start Fiber server: %v", err)
	}
}

func startCronJob(db *gorm.DB, browserPath string, pwHeadless bool, apiKey string, interval int) {
	pw, err := initializePlaywright()
	if err != nil {
		logrus.Fatalln("error during Playwright setup:", err)
	}
	defer func(pw *playwright.Playwright) {
		_ = pw.Stop()
	}(pw)

	c := syncronizer.Env{Browser: initializeBrowser(pw, browserPath, pwHeadless), Database: db, AlchemyApiKey: apiKey}
	runCronJob(c, interval)
}

func initializePlaywright() (*playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	if err := playwright.Install(); err != nil {
		return nil, err
	}

	return pw, nil
}

func initializeBrowser(pw *playwright.Playwright, browserPath string, headless bool) playwright.Browser {
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(browserPath),
		Headless:       playwright.Bool(headless),
	})
	if err != nil {
		logrus.Fatalln("error during browser launch:", err)
	}
	return browser
}

func runCronJob(c syncronizer.Env, interval int) {
	var mutex sync.Mutex
	s := gocron.NewScheduler(time.Local)

	for {
		_, err := s.Every(interval).Minutes().Do(func() {
			mutex.Lock()
			defer mutex.Unlock()
			c.Sync()
		})

		if err != nil {
			logrus.WithError(err).Fatalln("cannot start cron")
		}

		s.StartBlocking()
		logrus.Warning("scheduler stopped, shutting down...restarting")
	}
}
