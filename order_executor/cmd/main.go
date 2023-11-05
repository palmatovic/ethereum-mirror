package main

import (
	"fmt"
	"io"
	sync "order-executor/pkg/cron"
	db_t "order-executor/pkg/model/database/token"
	db_w "order-executor/pkg/model/database/wallet"
	db_wto "order-executor/pkg/model/database/wallet_token"
	db_wtr "order-executor/pkg/model/database/wallet_transaction"
	"os"
	"strings"
	syncronize "sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Environment
// Defines the structure for holding environment variables
type Environment struct {
	MinPercOrderThreshold        int     `env:"MIN_PERC_ORDER_THRESHOLD,required"`
	MaxPercOrderThreshold        int     `env:"MAX_PERC_ORDER_THRESHOLD,required"`
	MinAbsOrderThreshold         int     `env:"MIN_ABS_ORDER_THRESHOLD,required"`
	MaxAbsOrderThreshold         int     `env:"MAX_ABS_ORDER_THRESHOLD,require"`
	SetMaxPercThreshold          bool    `env:"SET_MAX_PERC_THRESHOLD,required"`
	SetMinPercThreshold          bool    `env:"SET_MIN_PERC_THRESHOLD,required"`
	OrderTimeExpirationThreshold int     `env:"ORDER_TIME_EXPIRATION_THRESHOLD,required"`
	StopEarningThreshold         int     `env:"STOP_EARNING_THRESHOLD,required"`
	StopLossThreshold            int     `env:"STOP_LOSS_THRESHOLD,required"`
	MaxPriceRangePerc            float32 `env:"MAX_PRICE_RANGE_PERC,required"`
	LogLevel                     string  `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath                  string  `env:"LOG_FILE_PATH" envDefault:"./orderexecutor.log"`
	ConsoleLogEnable             bool    `env:"CONSOLE_LOG_ENABLE" envDefault:"true"`
}

func main() {
	var (
		e   = Environment{}
		log = logrus.New()
		err error
	)

	// Parse environment variables into the 'e' struct
	if err = env.Parse(&e); err != nil {
		log.WithError(err).Fatalln("error during environment parsing")
	}

	// Open a connection to the SQLite database
	db := initializeDatabase()
	migrateDatabase(db)

	// Create an instance of the cron environment
	c := sync.Env{
		Database:                     db,
		MinPercOrderThreshold:        e.MinPercOrderThreshold,
		MaxPercOrderThreshold:        e.MaxPercOrderThreshold,
		MinAbsOrderThreshold:         e.MinAbsOrderThreshold,
		MaxAbsOrderThreshold:         e.MaxAbsOrderThreshold,
		SetMaxPercThreshold:          e.SetMaxPercThreshold,
		SetMinPercThreshold:          e.SetMinPercThreshold,
		StopEarningThreshold:         e.StopEarningThreshold,
		StopLossThreshold:            e.StopLossThreshold,
		OrderTimeExpirationThreshold: e.OrderTimeExpirationThreshold,
		MaxPriceRangePerc:            e.MaxPriceRangePerc,
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

		_, syncErr := c.ExecuteOrders()
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

func loadAppConfig() Environment {
	var config Environment
	if err := env.Parse(&config); err != nil {
		handleError(err, "error during environment parsing")
	}
	config.LogFilePath = fmt.Sprintf("%s_%s.log", strings.Split(config.LogFilePath, ".log")[0], time.Now().UTC().Format(time.RFC3339))
	return config
}

func initializeLogger(config Environment) {
	logrus.New()
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: false,
		PrettyPrint:       true,
	})

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	handleError(err, "error during parse log level")

	logrus.SetLevel(logLevel)

	logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err, "error during creation of log file")

	var multiWriter io.Writer
	if config.ConsoleLogEnable {
		multiWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		multiWriter = io.MultiWriter(logFile)
	}
	logrus.SetOutput(multiWriter)

}

func initializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./orderexecutor.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	handleError(err, "error during database connection")

	return db
}

func migrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&db_w.Wallet{}, &db_t.Token{}, &db_wto.WalletToken{}, &db_wtr.WalletTransaction{}, &db_t.TokenPrice{})
	handleError(err, "error during migration of database")
}

func handleError(err error, message string) {
	if err != nil {
		logrus.Fatalf("%s: %v", message, err)
	}
}
