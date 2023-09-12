package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	syncronize "sync"
	"time"
	"wallet-syncronizer/pkg/cron/sync"
	"wallet-syncronizer/pkg/database/token"
	"wallet-syncronizer/pkg/database/wallet"
	"wallet-syncronizer/pkg/database/wallet_token"
)

// Environment
// Defines the structure for holding environment variables
type Environment struct {
	Wallets            []string `env:"WALLETS,required"` // 0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367
	PlaywrightHeadLess bool     `env:"PLAYWRIGHT_HEADLESS,required"`
	AlchemyApiKey      string   `env:"ALCHEMY_API_KEY,required"` // owUCVigVvnHA63o0C6mh3yrf3jxMkV7b
}

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
	err = db.AutoMigrate(&wallet.Wallet{}, &token.Token{}, &wallet_token.WalletToken{})
	if err != nil {
		logrus.WithError(err).Fatalln("error during migration of database")
	}

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

	//executablePath := "/usr/bin/brave-browser"
	// Launch Firefox browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		//ExecutablePath: &executablePath,
		Headless: playwright.Bool(e.PlaywrightHeadLess),
	})
	if err != nil {
		logrus.Fatalln("error during browser launch:", err)
	}

	// Create an instance of the cron environment
	c := sync.Env{Browser: browser, Database: db, Wallets: e.Wallets, AlchemyApiKey: e.AlchemyApiKey}

	// Create a new cron scheduler
	s := gocron.NewScheduler(time.Local)

	// Create a mutex for synchronization
	var mutex syncronize.Mutex

	// Define the cron job using cron syntax
	_, err = s.Every(10).Minute().Do(func() {
		// Lock the mutex before starting the task
		mutex.Lock()
		defer mutex.Unlock() // Unlock the mutex when the function finishes
		c.Sync()
	})
	// Start the cron scheduler (blocking call)
	s.StartBlocking()

	// This point is reached after the scheduler stops (due to blocking nature)
	logrus.Infoln("scheduler stopped, shutting down")
}

// token sniffer https://gopluslabs.io/token-security/1/:contractaddress
