package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	syncronize "sync"
	"time"
	"transaction-extractor/pkg/cron/sync"
	"transaction-extractor/pkg/database/address"
	"transaction-extractor/pkg/database/address_status"
)

// Environment
// Defines the structure for holding environment variables
type Environment struct {
	Addresses          []string `env:"ADDRESSES,required"` // 0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367
	PlaywrightHeadLess bool     `env:"PLAYWRIGHT_HEADLESS,required"`
	AlchemyApiKey      string   `env:"ALCHEMY_API_KEY,required"` // owUCVigVvnHA63o0C6mh3yrf3jxMkV7b
}

func main() {
	var (
		e   = Environment{}
		log = logrus.New()
		db  *gorm.DB
		err error
	)

	// Parse environment variables into the 'e' struct
	if err = env.Parse(&e); err != nil {
		log.WithError(err).Fatalln("error during environment parsing")
	}

	// Open a connection to the SQLite database
	if db, err = gorm.Open(sqlite.Open("wallet.db"), &gorm.Config{}); err != nil {
		log.WithError(err).Fatalln("error during database connection")
	}

	// Perform automatic database schema migration
	err = db.AutoMigrate(&address.Address{}, &address_status.AddressStatus{})
	if err != nil {
		log.WithError(err).Fatalln("error during migration of database")
	}

	// Set up Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalln("error during Playwright startup:", err)
	}
	defer func(pw *playwright.Playwright) {
		_ = pw.Stop()
	}(pw)

	// Install Playwright
	if err = playwright.Install(); err != nil {
		log.Fatalln("error during Playwright installation:", err)
	}

	executablePath := "/usr/bin/brave-browser"
	// Launch Firefox browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: &executablePath,
		Headless:       playwright.Bool(e.PlaywrightHeadLess),
	})
	if err != nil {
		log.Fatalln("error during browser launch:", err)
	}

	// Create an instance of the cron environment
	c := sync.Env{Browser: browser, Database: db, Addresses: e.Addresses, AlchemyApiKey: e.AlchemyApiKey}

	// Create a new cron scheduler
	s := gocron.NewScheduler(time.Local)

	// Create a mutex for synchronization
	var mutex syncronize.Mutex

	// Define the cron job using cron syntax
	_, err = s.Every(1).Minute().Do(func() {
		// Lock the mutex before starting the task
		mutex.Lock()
		defer mutex.Unlock() // Unlock the mutex when the function finishes

		_, syncErr := c.SyncTransactions()
		if syncErr != nil {
			log.Errorln("error during database sync:", syncErr)
		} else {
			log.Infoln("database sync completed successfully")
		}
	})
	// Start the cron scheduler (blocking call)
	s.StartBlocking()

	// This point is reached after the scheduler stops (due to blocking nature)
	log.Infoln("scheduler stopped, shutting down")
}
