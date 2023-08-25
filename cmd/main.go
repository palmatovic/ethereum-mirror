package main

import (
	"context"
	"ethereum-mirror/pkg/cron"
	"ethereum-mirror/pkg/database"
	"github.com/playwright-community/playwright-go"
	scheduler "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger
	log := logrus.New()

	db, err := gorm.Open(sqlite.Open("scraping.db"), &gorm.Config{})

	err = db.AutoMigrate(&database.Transaction{}, &database.Scraping{})
	if err != nil {
		log.WithError(err).Fatalln("error during migration of database")
	}

	// Create a context that can be used for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Launch Firefox browser
	browser, err := pw.Firefox.Launch()
	if err != nil {
		log.Fatalln("error during browser launch:", err)
	}

	c := cron.Env{Browser: browser, Database: db}

	// Create a new cron scheduler
	cronScheduler := scheduler.New()

	// Define the cron job to run c.SyncTransactions every 1 minute
	_, err = cronScheduler.AddFunc("*/1 * * * *", func() {
		_, syncErr := c.SyncTransactions()
		if syncErr != nil {
			log.Errorln("error during database sync:", syncErr)
		} else {
			log.Infoln("database sync completed successfully")
		}
	})
	if err != nil {
		log.Fatalln("error scheduling cron job:", err)
	}

	// Start the cron scheduler
	cronScheduler.Start()

	// Set up a signal handler to gracefully shut down the program
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	select {
	case sig := <-signals:
		log.Infof("received signal %s, gracefully shutting down...", sig)
		cancel() // Trigger graceful shutdown

		// Wait for the cron jobs to finish before exiting
		<-time.After(time.Minute) // Adjust the wait duration as needed

		log.Infoln("shutdown complete")
	}
}
