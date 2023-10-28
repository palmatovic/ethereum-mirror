package main

import (
	company_api "auth/pkg/api/company"
	product_api "auth/pkg/api/product"
	"auth/pkg/cron/sync"
	"auth/pkg/database/company"
	"auth/pkg/database/product"
	"auth/pkg/init"
	token_util "auth/pkg/service_util/fiber/jwt/token"
	jwt_util "auth/pkg/service_util/fiber/jwt/validator"
	product_url "auth/pkg/url/product"
	"context"
	"crypto/rsa"
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

type appConfig struct {
	FiberPort              int    `env:"FIBER_PORT" envDefault:"3000"`
	SyncJobIntervalMinutes int    `env:"SYNC_JOB_INTERVAL_MINUTES" envDefault:"1"`
	LogLevel               string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFilePath            string `env:"LOG_FILE_PATH" envDefault:"./auth.log"`
	ConsoleLogEnable       bool   `env:"CONSOLE_LOG_ENABLE" envDefault:"true"`
	InitScriptFilepath     string `env:"INIT_SCRIPT_FILEPATH,required"`
}

func main() {
	config := loadAppConfig()

	initializeLogger(config)

	db := initializeDatabase()

	migrateDatabase(db)

	jwtPublicKey := init.NewService(db).Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var eg errgroup.Group

	eg.Go(func() error {
		return runSyncJob(ctx, db, config.SyncJobIntervalMinutes)
	})

	eg.Go(func() error {
		return runFiber(ctx, db, config.FiberPort, jwtPublicKey)
	})

	if err := eg.Wait(); err != nil {
		handleError(err, "application error")
	}
}

func loadAppConfig() appConfig {
	var config appConfig
	if err := env.Parse(&config); err != nil {
		handleError(err, "error during environment parsing")
	}
	config.LogFilePath = fmt.Sprintf("%s_%s.log", strings.Split(config.LogFilePath, ".log")[0], time.Now().UTC().Format(time.RFC3339))
	return config
}

func initializeLogger(config appConfig) {
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
		&product.Product{},
		&company.Company{},
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

func initializeFiberApp(db *gorm.DB, jwtPublicKey *rsa.PublicKey) *fiber.App {
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

	registerAPIRoutes(app, db, jwt_util.NewService(jwtPublicKey).Validator())

	return app
}

func registerAPIRoutes(app *fiber.App, db *gorm.DB, jwtValidator fiber.Handler) {
	//permApi := perm_api.NewApi(db)
	//permGraphQLApi := graphql_api.NewApi(perm_graphql_schema.Schema(db))

	productApi := product_api.NewApi(db)
	//productGraphQLApi := graphql_api.NewApi(product_graphql_schema.Schema(db))

	companyApi := company_api.NewApi(db)
	//companyGraphQLApi := graphql_api.NewApi(company_graphql_schema.Schema(db))

	//resourceApi := resource_api.NewApi(db)
	//resourceGraphQLApi := graphql_api.NewApi(resource_graphql_schema.Schema(db))
	//
	//resourcePermApi := resource_perm_api.NewApi(db)
	//resourcePermGraphQLApi := graphql_api.NewApi(resource_perm_graphql_schema.Schema(db))
	//
	//userApi := user_api.NewApi(db)
	//userGraphQLApi := graphql_api.NewApi(user_graphql_schema.Schema(db))
	//
	//userProductApi := user_product_api.NewApi(db)
	//userProductGraphQLApi := graphql_api.NewApi(user_product_graphql_schema.Schema(db))
	//
	//userResourceApi := user_resource_api.NewApi(db)
	//userResourceGraphQLApi := graphql_api.NewApi(user_resource_graphql_schema.Schema(db))
	//
	//userResourcePermApi := user_resource_perm_api.NewApi(db)
	//userResourcePermGraphQLApi := graphql_api.NewApi(user_resource_perm_graphql_schema.Schema(db))
	apiList := []struct {
		method  string
		path    string
		handler []fiber.Handler
	}{
		//{"POST", perm_url.GraphQL, permGraphQLApi.Post},
		//{"GET", perm_url.Get, permApi.Get},
		//{"GET", perm_url.List, permApi.List},
		//{"POST", perm_url.Create, permApi.Create},
		//{"PUT", perm_url.Update, permApi.Update},
		//{"DELETE", perm_url.Delete, permApi.Delete},

		//{"POST", product_url.GraphQL, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productGraphQLApi.Post}},
		{"GET", product_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Get}},
		{"GET", product_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.List}},
		{"POST", product_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Create}},
		{"PUT", product_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Update}},
		{"DELETE", product_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Delete}},

		//{"POST", product_url.GraphQL, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyGraphQLApi.Post}},
		{"GET", product_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Get}},
		{"GET", product_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.List}},
		{"POST", product_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Create}},
		{"PUT", product_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Update}},
		{"DELETE", product_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Delete}},

		//{"POST", resource_url.GraphQL, resourceGraphQLApi.Post},
		//{"GET", resource_url.Get, resourceApi.Get},
		//{"GET", resource_url.List, resourceApi.List},
		//{"POST", resource_url.Create, resourceApi.Create},
		//{"PUT", resource_url.Update, resourceApi.Update},
		//{"DELETE", resource_url.Delete, resourceApi.Delete},
		//
		//{"POST", resource_perm_url.GraphQL, resourcePermGraphQLApi.Post},
		//{"GET", resource_perm_url.Get, resourcePermApi.Get},
		//{"GET", resource_perm_url.List, resourcePermApi.List},
		//{"POST", resource_perm_url.Create, resourcePermApi.Create},
		//{"PUT", resource_perm_url.Update, resourcePermApi.Update},
		//{"DELETE", resource_perm_url.Delete, resourcePermApi.Delete},
		//
		//{"POST", user_url.GraphQL, userGraphQLApi.Post},
		//{"GET", user_url.Get, userApi.Get},
		//{"GET", user_url.List, userApi.List},
		//{"POST", user_url.Create, userApi.Create},
		//{"PUT", user_url.Update, userApi.Update},
		//{"DELETE", user_url.Delete, userApi.Delete},
		//
		//{"POST", user_product_url.GraphQL, userProductGraphQLApi.Post},
		//{"GET", user_product_url.Get, userProductApi.Get},
		//{"GET", user_product_url.List, userProductApi.List},
		//{"POST", user_product_url.Create, userProductApi.Create},
		//{"PUT", user_product_url.Update, userProductApi.Update},
		//{"DELETE", user_product_url.Delete, userProductApi.Delete},
		//
		//{"POST", user_resource_url.GraphQL, userResourceGraphQLApi.Post},
		//{"GET", user_resource_url.Get, userResourceApi.Get},
		//{"GET", user_resource_url.List, userResourceApi.List},
		//{"POST", user_resource_url.Create, userResourceApi.Create},
		//{"PUT", user_resource_url.Update, userResourceApi.Update},
		//{"DELETE", user_resource_url.Delete, userResourceApi.Delete},
		//
		//{"POST", user_resource_perm_url.GraphQL, userResourcePermGraphQLApi.Post},
		//{"GET", user_resource_perm_url.Get, userResourcePermApi.Get},
		//{"GET", user_resource_perm_url.List, userResourcePermApi.List},
		//{"POST", user_resource_perm_url.Create, userResourcePermApi.Create},
		//{"PUT", user_resource_perm_url.Update, userResourcePermApi.Update},
		//{"DELETE", user_resource_perm_url.Delete, userResourcePermApi.Delete},
	}

	for _, api := range apiList {
		app.Add(api.method, api.path, api.handler...)
	}
}

func runFiber(ctx context.Context, db *gorm.DB, port int, jwtPublicKey *rsa.PublicKey) error {
	app := initializeFiberApp(db, jwtPublicKey)
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
