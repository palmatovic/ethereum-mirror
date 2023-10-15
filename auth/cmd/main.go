package main

import (
	"auth/pkg/cron/sync"
	"auth/pkg/database/perm"
	"auth/pkg/database/product"
	"auth/pkg/database/resource"
	"auth/pkg/database/resource_perm"
	"auth/pkg/database/user"
	"auth/pkg/database/user_product"
	"auth/pkg/database/user_resource"
	"auth/pkg/database/user_resource_perm"
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	FiberPort              int    `env:"FIBER_PORT" envDefault:"3000"`
	SyncJobIntervalMinutes int    `env:"SYNC_JOB_INTERVAL_MINUTES" envDefault:"1"`
	LogLevel               string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath            string `env:"LOG_FILE_PATH" envDefault:"./auth.log"`
	ConsoleLogEnable       bool   `env:"CONSOLE_LOG_ENABLE" envDefault:"true"`
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
		return runSyncJob(ctx, db, config.SyncJobIntervalMinutes)
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
	db, err := gorm.Open(sqlite.Open("./auth.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	handleError(err, "error during database connection")

	return db
}

func migrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&perm.Perm{},
		&product.Product{},
		&resource.Resource{},
		&resource_perm.ResourcePerm{},
		&user.User{},
		&user_product.UserProduct{},
		&user_resource.UserResource{},
		&user_resource_perm.UserResourcePerm{},
	)
	handleError(err, "error during migration of database")
}

func runSyncJob(ctx context.Context, db *gorm.DB, interval int) error {
	c := sync.NewSync(db)
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
	//tokenApi := token_api.NewApi(db)
	//tokenGraphqlApi := graph_api.NewApi(token_graphql_schema.Schema(db))
	//
	//walletApi := wallet_api.NewApi(db)
	//walletGraphqlApi := graph_api.NewApi(wallet_graphql_schema.Schema(db))
	//
	//walletTokenApi := wallet_token_api.NewApi(db)
	//walletTokenGraphqlApi := graph_api.NewApi(wallet_token_graphql_schema.Schema(db))
	//
	//walletTransactionApi := wallet_transaction_api.NewApi(db)
	//walletTransactionGraphqlApi := graph_api.NewApi(wallet_transaction_graphql_schema.Schema(db))
	//
	apiList := []struct {
		method  string
		path    string
		handler fiber.Handler
	}{
		//{"POST", token_url.GraphQL, tokenGraphqlApi.Post},
		//{"GET", token_url.Get, tokenApi.Get},
		//{"GET", token_url.List, tokenApi.List},
		//{"POST", wallet_url.GraphQL, walletGraphqlApi.Post},
		//{"GET", wallet_url.Get, walletApi.Get},
		//{"GET", wallet_url.List, walletApi.List},
		//{"POST", wallet_url.Create, walletApi.Create},
		////{"PUT", wallet_url.Update, walletApi.Update},
		//{"DELETE", wallet_url.Delete, walletApi.Delete},
		//{"POST", wallet_token_url.GraphQL, walletTokenGraphqlApi.Post},
		//{"GET", wallet_token_url.Get, walletTokenApi.Get},
		//{"GET", wallet_token_url.List, walletTokenApi.List},
		//{"POST", wallet_transaction_url.GraphQL, walletTransactionGraphqlApi.Post},
		//{"GET", wallet_transaction_url.Get, walletTransactionApi.Get},
		//{"GET", wallet_transaction_url.List, walletTransactionApi.List},
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
			return nil
		default:
			handleError(err, "cannot start Fiber server")
		}
	}
	return nil
}
