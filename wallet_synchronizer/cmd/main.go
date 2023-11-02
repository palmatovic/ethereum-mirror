package main

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"os"
	"time"
	graphql_api "wallet-synchronizer/pkg/api/graphql"
	token_api "wallet-synchronizer/pkg/api/token"
	"wallet-synchronizer/pkg/database/wallet"
	"wallet-synchronizer/pkg/database/wallet_token"
	"wallet-synchronizer/pkg/database/wallet_transaction"
	token_graphql_schema "wallet-synchronizer/pkg/graphql/schema/token"
	wallet_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet"
	wallet_token_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet_token"
	wallet_transaction_graphql_schema "wallet-synchronizer/pkg/graphql/schema/wallet_transaction"
	init_database "wallet-synchronizer/pkg/init/database"
	"wallet-synchronizer/pkg/init/environment"
	init_fiber "wallet-synchronizer/pkg/init/fiber"
	init_logger "wallet-synchronizer/pkg/init/logger"
	init_playwright "wallet-synchronizer/pkg/init/playwright"
	"wallet-synchronizer/pkg/perm"
	"wallet-synchronizer/pkg/resource"
	"wallet-synchronizer/pkg/service_util/aes"
	"wallet-synchronizer/pkg/service_util/fiber/jwt/token"
	"wallet-synchronizer/pkg/service_util/fiber/jwt/validator"
	"wallet-synchronizer/pkg/service_util/rsa"
	token_url "wallet-synchronizer/pkg/url/token"
	wallet_url "wallet-synchronizer/pkg/url/wallet"
	wallet_token_url "wallet-synchronizer/pkg/url/wallet_token"
	wallet_transaction_url "wallet-synchronizer/pkg/url/wallet_transaction"

	wallet_api "wallet-synchronizer/pkg/api/wallet"
	wallet_token_api "wallet-synchronizer/pkg/api/wallet_token"
	wallet_transaction_api "wallet-synchronizer/pkg/api/wallet_transaction"
	syncronizer "wallet-synchronizer/pkg/cron/sync"
)

func main() {
	config, err := environment.NewService().Init()
	handleError(err, "init environment error")

	err = init_logger.NewService(config.LogLevel, config.LogFilePath, config.ConsoleLogEnabled).Init()
	handleError(err, "init logger error")

	aes256EncryptionKey := aes.Key(config.AES256EncryptionKey)

	db, err := init_database.NewService(
		&aes256EncryptionKey,
		"./wallet_synchronizer.db",
		config.OwnWallet,
		&wallet.Wallet{},
		&wallet_token.WalletToken{},
		&wallet_transaction.WalletTransaction{},
	).Init()
	handleError(err, "init database error")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	browser, err := init_playwright.NewService(config.PlaywrightHeadless).Init()
	handleError(err, "init playwright error")

	var eg errgroup.Group

	eg.Go(func() error {
		return runSyncJob(ctx, db, *browser, config.OwnWallet, config.AlchemyAPIKey, config.ScrapeIntervalMinutes)
	})

	publicKey, err := rsa.PublicKeyFilepath(config.RSA256PublicKeyFilepath).ConvertToObj()
	handleError(err, "get auth public key error")

	serverSslCert, err := os.ReadFile(config.ServerSSLCertFilepath)
	handleError(err, "load server ssl cert error")

	serverSslKey, err := os.ReadFile(config.ServerSSLKeyFilepath)
	handleError(err, "load server ssl key error")

	sslCert, err := tls.X509KeyPair(serverSslCert, serverSslKey)
	handleError(err, "get server ssl error")

	jwtValidator := validator.NewService(publicKey).Validator()

	tokenApi := token_api.NewApi(db)
	tokenGraphqlApi := graphql_api.NewApi(token_graphql_schema.Schema(db))

	walletApi := wallet_api.NewApi(db)
	walletGraphqlApi := graphql_api.NewApi(wallet_graphql_schema.Schema(db))

	walletTokenApi := wallet_token_api.NewApi(db)
	walletTokenGraphqlApi := graphql_api.NewApi(wallet_token_graphql_schema.Schema(db))

	walletTransactionApi := wallet_transaction_api.NewApi(db)
	walletTransactionGraphqlApi := graphql_api.NewApi(wallet_transaction_graphql_schema.Schema(db))

	eg.Go(func() error {
		return init_fiber.NewService(sslCert, ctx, config.FiberPort, []init_fiber.Api{
			init_fiber.NewApi("POST", token_url.GraphQL, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Token, perm.GraphQL), tokenGraphqlApi.Post}),
			init_fiber.NewApi("GET", token_url.Get, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Token, perm.Get), tokenApi.Get}),
			init_fiber.NewApi("GET", token_url.List, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Token, perm.List), tokenApi.List}),

			init_fiber.NewApi("POST", wallet_url.GraphQL, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Wallet, perm.GraphQL), walletGraphqlApi.Post}),
			init_fiber.NewApi("GET", wallet_url.Get, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Wallet, perm.Get), walletApi.Get}),
			init_fiber.NewApi("GET", wallet_url.List, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Wallet, perm.List), walletApi.List}),
			init_fiber.NewApi("POST", wallet_url.Create, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Wallet, perm.Create), walletApi.Create}),
			init_fiber.NewApi("DELETE", wallet_url.Delete, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.Wallet, perm.Delete), walletApi.Delete}),

			init_fiber.NewApi("POST", wallet_token_url.GraphQL, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletToken, perm.GraphQL), walletTokenGraphqlApi.Post}),
			init_fiber.NewApi("GET", wallet_token_url.Get, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletToken, perm.Get), walletTokenApi.Get}),
			init_fiber.NewApi("GET", wallet_token_url.List, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletToken, perm.List), walletTokenApi.List}),

			init_fiber.NewApi("POST", wallet_transaction_url.GraphQL, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletTransaction, perm.GraphQL), walletTransactionGraphqlApi.Post}),
			init_fiber.NewApi("GET", wallet_transaction_url.Get, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletTransaction, perm.Get), walletTransactionApi.Get}),
			init_fiber.NewApi("GET", wallet_transaction_url.List, []fiber.Handler{jwtValidator, token.AccessTokenValidator, token.HasResourcePerm(resource.WalletTransaction, perm.List), walletTransactionApi.List}),
		}).Init()
	})

	if err = eg.Wait(); err != nil {
		handleError(err, "application error")
	}

}

func runSyncJob(ctx context.Context, db *gorm.DB, browser playwright.Browser, ownWallet string, apiKey string, interval int64) error {

	var ownWalletDb wallet.Wallet
	if err := db.Where("WalletId = ?", ownWallet).First(&ownWalletDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(err, "own wallet not found")
		}
	}

	c := syncronizer.NewSync(browser, db, ownWalletDb, apiKey)

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
