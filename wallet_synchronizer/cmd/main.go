package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
	"time"
	graph_api "wallet-synchronizer/pkg/api/graphql"
	token_api "wallet-synchronizer/pkg/api/token"
	token_graphql_schema "wallet-synchronizer/pkg/graphql/schema/token"
	wallet_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet"
	wallet_token_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet_token"
	wallet_transaction_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet_transaction"

	wallet_api "wallet-synchronizer/pkg/api/wallet"
	wallet_token_api "wallet-synchronizer/pkg/api/wallet_token"
	wallet_transaction_api "wallet-synchronizer/pkg/api/wallet_transaction"
	syncronizer "wallet-synchronizer/pkg/cron/sync"
	"wallet-synchronizer/pkg/database/token"
	"wallet-synchronizer/pkg/database/wallet"
	"wallet-synchronizer/pkg/database/wallet_token"
	"wallet-synchronizer/pkg/database/wallet_transaction"
	token_url "wallet-synchronizer/pkg/util/url/token"
	wallet_url "wallet-synchronizer/pkg/util/url/wallet"
	wallet_token_url "wallet-synchronizer/pkg/util/url/wallet_token"
	wallet_transaction_url "wallet-synchronizer/pkg/util/url/wallet_transaction"
)

type AppConfig struct {
	PlaywrightHeadless    bool   `env:"PLAYWRIGHT_HEADLESS" envDefault:"true"`
	AlchemyAPIKey         string `env:"ALCHEMY_API_KEY" envDefault:"owUCVigVvnHA63o0C6mh3yrf3jxMkV7b"`
	FiberPort             int    `env:"FIBER_PORT" envDefault:"3000"`
	BrowserPath           string `env:"BROWSER_PATH" envDefault:"/usr/bin/brave-browser"`
	ScrapeIntervalMinutes int    `env:"SCRAPE_INTERVAL_MINUTES" envDefault:"1"`
	LogLevel              string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath           string `env:"LOG_FILE_PATH" envDefault:"./wallet_synchronizer.log"`
}

func main() {
	config := loadAppConfig()
	initializeLogger(config)

	db := initializeDatabase()
	migrateDatabase(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var eg errgroup.Group

	eg.Go(func() error {
		return runSyncJob(ctx, db, config.BrowserPath, config.PlaywrightHeadless, config.AlchemyAPIKey, config.ScrapeIntervalMinutes)
	})

	eg.Go(func() error {
		return runFiber(ctx, db, config.FiberPort)
	})

	if err := eg.Wait(); err != nil {
		handleError(err, "application error")
	}
}

func loadAppConfig() AppConfig {
	var config AppConfig
	if err := env.Parse(&config); err != nil {
		handleError(err, "error during environment parsing")
	}
	config.LogFilePath = fmt.Sprintf("%s_%s.log", strings.Split(config.LogFilePath, ".log")[0], time.Now().UTC().Format(time.RFC3339))
	return config
}

func initializeLogger(config AppConfig) {
	logrus.New()
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
	})

	// Set log level based on the configuration
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel // Default to Info if the provided log level is invalid
	}
	logrus.SetLevel(logLevel)

	// If LogFilePath is provided, log to a file
	if config.LogFilePath != "" {
		logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(logFile)
		} else {
			logrus.Warn("Failed to log to file, using default stderr")
		}
	}
}

func initializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./wallet_synchronizer.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	handleError(err, "error during database connection")

	return db
}

func migrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&wallet.Wallet{}, &token.Token{}, &wallet_token.WalletToken{}, &wallet_transaction.WalletTransaction{})
	handleError(err, "error during migration of database")
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

func initializePlaywright() (*playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	if err = playwright.Install(&playwright.RunOptions{Verbose: false}); err != nil {
		return nil, err
	}

	return pw, nil
}

func initializeBrowser(pw *playwright.Playwright, browserPath string, headless bool) playwright.Browser {
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:       playwright.Bool(headless),
		ExecutablePath: playwright.String(browserPath),
	})
	handleError(err, "error during browser launch")
	return browser
}

func runSyncJob(ctx context.Context, db *gorm.DB, browserPath string, headless bool, apiKey string, interval int) error {
	pw, err := initializePlaywright()
	if err != nil {
		return err
	}
	defer func() {
		if err := pw.Stop(); err != nil {
			handleError(err, "error stopping Playwright")
		}
	}()

	c := syncronizer.NewSync(initializeBrowser(pw, browserPath, headless), db, apiKey)

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			c.Sync()
		}
	}
}

func handleError(err error, message string) {
	if err != nil {
		logrus.Fatalf("%s: %v", message, err)
	}
}

func initializeFiberApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
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

	registerAPIRoutes(app, db)

	return app
}

func registerAPIRoutes(app *fiber.App, db *gorm.DB) {
	tokenApi := token_api.NewApi(db)
	tokenGraphqlApi := graph_api.NewApi(token_graphql_schema.Schema(db))

	walletApi := wallet_api.NewApi(db)
	walletGraphqlApi := graph_api.NewApi(wallet_graphql_schema.Schema(db))

	walletTokenApi := wallet_token_api.NewApi(db)
	walletTokenGraphqlApi := graph_api.NewApi(wallet_token_graphql_schema.Schema(db))

	walletTransactionApi := wallet_transaction_api.NewApi(db)
	walletTransactionGraphqlApi := graph_api.NewApi(wallet_transaction_graphql_schema.Schema(db))

	apiList := []struct {
		method  string
		path    string
		handler fiber.Handler
	}{
		{"POST", token_url.GraphQL, tokenGraphqlApi.Post},
		{"GET", token_url.Get, tokenApi.Get},
		{"GET", token_url.List, tokenApi.List},
		{"POST", wallet_url.GraphQL, walletGraphqlApi.Post},
		{"GET", wallet_url.Get, walletApi.Get},
		{"GET", wallet_url.List, walletApi.List},
		{"POST", wallet_url.Create, walletApi.Create},
		//{"PUT", wallet_url.Update, walletApi.Update},
		{"DELETE", wallet_url.Delete, walletApi.Delete},
		{"POST", wallet_token_url.GraphQL, walletTokenGraphqlApi.Post},
		{"GET", wallet_token_url.Get, walletTokenApi.Get},
		{"GET", wallet_token_url.List, walletTokenApi.List},
		{"POST", wallet_transaction_url.GraphQL, walletTransactionGraphqlApi.Post},
		{"GET", wallet_transaction_url.Get, walletTransactionApi.Get},
		{"GET", wallet_transaction_url.List, walletTransactionApi.List},
	}

	for _, api := range apiList {
		app.Add(api.method, api.path, api.handler)
	}
}

func runFiber(ctx context.Context, db *gorm.DB, port int) error {
	app := initializeFiberApp(db)
	app.Server()

	addr := fmt.Sprintf(":%d", port)
	err := app.Listen(addr)
	if err != nil {
		select {
		case <-ctx.Done():
			// Application is shutting down, don't treat this as an error.
			return nil
		default:
			// A real error occurred.
			handleError(err, "cannot start Fiber server")
			return err
		}
	}
	return nil
}
