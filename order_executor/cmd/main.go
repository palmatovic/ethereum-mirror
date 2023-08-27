package main

import (
	sync "order-executor/pkg/cron"
	database "order-executor/pkg/model/database"
	syncronize "sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Environment
// Defines the structure for holding environment variables
type Environment struct {
	MinPercOrderThreshold int  `env:"MINPERCORDERTHRESHOLD,required"` // 0x905615DE62BE9B1a6582843E8ceDeDB6BDA42367
	MaxPercOrderThreshold int  `env:"MAXPERCORDERTHRESHOLD,required"`
	SetMaxPercThreshold   bool `env:"SETMAXPERCTHRESHOLD,required"`
	SetMinPercThreshold   bool `env:"SETMINPERCTHRESHOLD,required"`
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
	if db, err = gorm.Open(sqlite.Open("orderexecutor.db"), &gorm.Config{}); err != nil {
		log.WithError(err).Fatalln("error during database connection")
	}

	// Perform automatic database schema migration
	err = db.AutoMigrate(&database.Movement{}, &database.Transaction{}, &database.Order{})
	if err != nil {
		log.WithError(err).Fatalln("error during migration of database")
	}

	// Create an instance of the cron environment
	c := sync.Env{
		Database:              db,
		MinPercOrderThreshold: e.MinPercOrderThreshold,
		MaxPercOrderThreshold: e.MaxPercOrderThreshold,
		SetMaxPercThreshold:   e.SetMaxPercThreshold,
		SetMinPercThreshold:   e.SetMinPercThreshold,
	}

	// Create a new cron scheduler
	s := gocron.NewScheduler(time.Local)

	// Create a mutex for synchronization
	var mutex syncronize.Mutex

	// Define the cron job using cron syntax
	_, err = s.Every(1).Minute().Do(func() {
		// Lock the mutex before starting the task
		mutex.Lock()
		defer mutex.Unlock() // Unlock the mutex when the function finishes

		_, syncErr := c.ExecuteOrdres()
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
