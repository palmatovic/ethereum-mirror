package main

import (
	"auth/pkg/api/company"
	"auth/pkg/api/group"
	"auth/pkg/api/group_role"
	"auth/pkg/api/group_role_resource_perm"
	"auth/pkg/api/product"
	"auth/pkg/cron/sync"
	company_db "auth/pkg/database/company"
	group_db "auth/pkg/database/group"
	group_role_db "auth/pkg/database/group_role"
	group_role_resource_perm_db "auth/pkg/database/group_role_resource_perm"
	perm_db "auth/pkg/database/perm"
	product_db "auth/pkg/database/product"
	resource_db "auth/pkg/database/resource"
	resource_perm_db "auth/pkg/database/resource_perm"
	role_db "auth/pkg/database/role"
	user_db "auth/pkg/database/user"
	user_group_role_db "auth/pkg/database/user_group_role"
	user_product_db "auth/pkg/database/user_product"
	init_database "auth/pkg/init/database"
	"auth/pkg/init/environment"
	init_fiber "auth/pkg/init/fiber"
	init_logger "auth/pkg/init/logger"
	"auth/pkg/service/product/get_public_key"
	token_util "auth/pkg/service_util/fiber/jwt/token"
	"auth/pkg/service_util/fiber/jwt/validator"
	company_url "auth/pkg/url/company"
	group_url "auth/pkg/url/group"
	group_role_url "auth/pkg/url/group_role"
	group_role_resource_perm_url "auth/pkg/url/group_role_resource_perm"
	product_url "auth/pkg/url/product"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

func main() {
	config, err := environment.NewService().Init()
	handleError(err, "init environment error")

	err = init_logger.NewService(config.LogLevel, config.LogFilePath, config.ConsoleLogEnabled).Init()
	handleError(err, "init logger error")

	db, err := init_database.NewService("./auth.db",
		&company_db.Company{},
		&group_db.Group{},
		&group_role_db.GroupRole{},
		&group_role_resource_perm_db.GroupRoleResourcePerm{},
		&perm_db.Perm{},
		&product_db.Product{},
		&resource_db.Resource{},
		&resource_perm_db.ResourcePerm{},
		&role_db.Role{},
		&user_db.User{},
		&user_group_role_db.UserGroupRole{},
		&user_product_db.UserProduct{},
	).Init()
	handleError(err, "init database error")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var eg errgroup.Group

	eg.Go(func() error {
		return runSyncJob(ctx, db, config.SyncJobIntervalMinutes)
	})

	_, publicKey, err := get_public_key.NewService(db, "auth").Get()
	handleError(err, "get auth public key error")

	jwtValidator := validator.NewService(publicKey).Validator()

	productApi := product.NewApi(db)
	companyApi := company.NewApi(db)
	groupApi := group.NewApi(db)
	groupRoleApi := group_role.NewApi(db)
	groupRoleResourcePermApi := group_role_resource_perm.NewApi(db)

	eg.Go(func() error {
		return init_fiber.NewService(ctx, config.FiberPort, []init_fiber.Api{
			init_fiber.NewApi("GET", product_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Get}),
			init_fiber.NewApi("GET", product_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.List}),
			init_fiber.NewApi("POST", product_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Create}),
			init_fiber.NewApi("PUT", product_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Update}),
			init_fiber.NewApi("DELETE", product_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, productApi.Delete}),

			init_fiber.NewApi("GET", company_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Get}),
			init_fiber.NewApi("GET", company_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.List}),
			init_fiber.NewApi("POST", company_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Create}),
			init_fiber.NewApi("PUT", company_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Update}),
			init_fiber.NewApi("DELETE", company_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, companyApi.Delete}),

			init_fiber.NewApi("GET", group_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupApi.Get}),
			init_fiber.NewApi("GET", group_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupApi.List}),
			init_fiber.NewApi("POST", group_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupApi.Create}),
			init_fiber.NewApi("PUT", group_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupApi.Update}),
			init_fiber.NewApi("DELETE", group_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupApi.Delete}),

			init_fiber.NewApi("GET", group_role_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleApi.Get}),
			init_fiber.NewApi("GET", group_role_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleApi.List}),
			init_fiber.NewApi("POST", group_role_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleApi.Create}),
			init_fiber.NewApi("PUT", group_role_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleApi.Update}),
			init_fiber.NewApi("DELETE", group_role_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleApi.Delete}),

			init_fiber.NewApi("GET", group_role_resource_perm_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleResourcePermApi.Get}),
			init_fiber.NewApi("GET", group_role_resource_perm_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleResourcePermApi.List}),
			init_fiber.NewApi("POST", group_role_resource_perm_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleResourcePermApi.Create}),
			init_fiber.NewApi("PUT", group_role_resource_perm_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleResourcePermApi.Update}),
			init_fiber.NewApi("DELETE", group_role_resource_perm_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, groupRoleResourcePermApi.Delete}),
		}).Init()
	})

	if err = eg.Wait(); err != nil {
		handleError(err, "application error")
	}
}

func runSyncJob(ctx context.Context, db *gorm.DB, interval int64) error {
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
