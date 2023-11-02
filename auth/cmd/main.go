package main

import (
	"auth/pkg/api/company"
	"auth/pkg/api/group"
	"auth/pkg/api/group_role"
	"auth/pkg/api/group_role_resource_perm"
	"auth/pkg/api/perm"
	"auth/pkg/api/product"
	"auth/pkg/api/resource"
	"auth/pkg/api/resource_perm"
	"auth/pkg/api/role"
	"auth/pkg/api/user"
	"auth/pkg/api/user_group_role"
	"auth/pkg/api/user_product"
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
	perm_constants "auth/pkg/perm"
	resource_constants "auth/pkg/resource"
	"auth/pkg/service/product/get_public_key"
	"auth/pkg/service/product/get_server_ssl"
	"auth/pkg/service_util/aes"
	token_util "auth/pkg/service_util/fiber/jwt/token"
	"auth/pkg/service_util/fiber/jwt/validator"
	company_url "auth/pkg/url/company"
	group_url "auth/pkg/url/group"
	group_role_url "auth/pkg/url/group_role"
	group_role_resource_perm_url "auth/pkg/url/group_role_resource_perm"
	perm_url "auth/pkg/url/perm"
	product_url "auth/pkg/url/product"
	resource_url "auth/pkg/url/resource"
	resource_perm_url "auth/pkg/url/resource_perm"
	role_url "auth/pkg/url/role"
	user_url "auth/pkg/url/user"
	user_group_role_url "auth/pkg/url/user_group_role"
	user_product_url "auth/pkg/url/user_product"
	"context"
	"crypto/tls"
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

	aes256EncryptionKey := aes.Key(config.AES256EncryptionKey)

	db, err := init_database.NewService(&aes256EncryptionKey, "./auth.db",
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

	_, publicKey, err := get_public_key.NewService(db, &aes256EncryptionKey, "auth").Get()
	handleError(err, "get auth public key error")

	_, serverSslCert, serverSslKey, err := get_server_ssl.NewService(db, &aes256EncryptionKey, "auth").Get()
	handleError(err, "get auth server ssl error")

	sslCert, err := tls.X509KeyPair(serverSslCert, serverSslKey)
	handleError(err, "get server ssl error")

	jwtValidator := validator.NewService(publicKey).Validator()

	productApi := product.NewApi(db, &aes256EncryptionKey)
	companyApi := company.NewApi(db)
	groupApi := group.NewApi(db)
	groupRoleApi := group_role.NewApi(db)
	groupRoleResourcePermApi := group_role_resource_perm.NewApi(db)
	permApi := perm.NewApi(db)
	resourceApi := resource.NewApi(db)
	resourcePermApi := resource_perm.NewApi(db)
	roleApi := role.NewApi(db)
	userApi := user.NewApi(db)
	userGroupRoleApi := user_group_role.NewApi(db)
	userProductApi := user_product.NewApi(db)

	eg.Go(func() error {
		return init_fiber.NewService(sslCert, ctx, config.FiberPort, []init_fiber.Api{

			init_fiber.NewApi("GET", product_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Product, perm_constants.Get), productApi.Get}),
			init_fiber.NewApi("GET", product_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Product, perm_constants.List), productApi.List}),
			init_fiber.NewApi("POST", product_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Product, perm_constants.Create), productApi.Create}),
			init_fiber.NewApi("PUT", product_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Product, perm_constants.Update), productApi.Update}),
			init_fiber.NewApi("DELETE", product_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Product, perm_constants.Delete), productApi.Delete}),

			init_fiber.NewApi("GET", company_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Company, perm_constants.Get), companyApi.Get}),
			init_fiber.NewApi("GET", company_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Company, perm_constants.List), companyApi.List}),
			init_fiber.NewApi("POST", company_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Company, perm_constants.Create), companyApi.Create}),
			init_fiber.NewApi("PUT", company_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Company, perm_constants.Update), companyApi.Update}),
			init_fiber.NewApi("DELETE", company_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Company, perm_constants.Delete), companyApi.Delete}),

			init_fiber.NewApi("GET", group_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Group, perm_constants.Get), groupApi.Get}),
			init_fiber.NewApi("GET", group_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Group, perm_constants.List), groupApi.List}),
			init_fiber.NewApi("POST", group_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Group, perm_constants.Create), groupApi.Create}),
			init_fiber.NewApi("PUT", group_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Group, perm_constants.Update), groupApi.Update}),
			init_fiber.NewApi("DELETE", group_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Group, perm_constants.Delete), groupApi.Delete}),

			init_fiber.NewApi("GET", group_role_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRole, perm_constants.Get), groupRoleApi.Get}),
			init_fiber.NewApi("GET", group_role_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRole, perm_constants.List), groupRoleApi.List}),
			init_fiber.NewApi("POST", group_role_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRole, perm_constants.Create), groupRoleApi.Create}),
			init_fiber.NewApi("PUT", group_role_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRole, perm_constants.Update), groupRoleApi.Update}),
			init_fiber.NewApi("DELETE", group_role_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRole, perm_constants.Delete), groupRoleApi.Delete}),

			init_fiber.NewApi("GET", group_role_resource_perm_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRoleResourcePerm, perm_constants.Get), groupRoleResourcePermApi.Get}),
			init_fiber.NewApi("GET", group_role_resource_perm_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRoleResourcePerm, perm_constants.List), groupRoleResourcePermApi.List}),
			init_fiber.NewApi("POST", group_role_resource_perm_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRoleResourcePerm, perm_constants.Create), groupRoleResourcePermApi.Create}),
			init_fiber.NewApi("PUT", group_role_resource_perm_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRoleResourcePerm, perm_constants.Update), groupRoleResourcePermApi.Update}),
			init_fiber.NewApi("DELETE", group_role_resource_perm_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.GroupRoleResourcePerm, perm_constants.Delete), groupRoleResourcePermApi.Delete}),

			init_fiber.NewApi("GET", perm_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Perm, perm_constants.Get), permApi.Get}),
			init_fiber.NewApi("GET", perm_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Perm, perm_constants.List), permApi.List}),
			init_fiber.NewApi("POST", perm_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Perm, perm_constants.Create), permApi.Create}),
			init_fiber.NewApi("PUT", perm_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Perm, perm_constants.Update), permApi.Update}),
			init_fiber.NewApi("DELETE", perm_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Perm, perm_constants.Delete), permApi.Delete}),

			init_fiber.NewApi("GET", resource_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Resource, perm_constants.Get), resourceApi.Get}),
			init_fiber.NewApi("GET", resource_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Resource, perm_constants.List), resourceApi.List}),
			init_fiber.NewApi("POST", resource_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Resource, perm_constants.Create), resourceApi.Create}),
			init_fiber.NewApi("PUT", resource_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Resource, perm_constants.Update), resourceApi.Update}),
			init_fiber.NewApi("DELETE", resource_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Resource, perm_constants.Delete), resourceApi.Delete}),

			init_fiber.NewApi("GET", resource_perm_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.ResourcePerm, perm_constants.Get), resourcePermApi.Get}),
			init_fiber.NewApi("GET", resource_perm_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.ResourcePerm, perm_constants.List), resourcePermApi.List}),
			init_fiber.NewApi("POST", resource_perm_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.ResourcePerm, perm_constants.Create), resourcePermApi.Create}),
			init_fiber.NewApi("PUT", resource_perm_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.ResourcePerm, perm_constants.Update), resourcePermApi.Update}),
			init_fiber.NewApi("DELETE", resource_perm_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.ResourcePerm, perm_constants.Delete), resourcePermApi.Delete}),

			init_fiber.NewApi("GET", role_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Role, perm_constants.Get), roleApi.Get}),
			init_fiber.NewApi("GET", role_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Role, perm_constants.List), roleApi.List}),
			init_fiber.NewApi("POST", role_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Role, perm_constants.Create), roleApi.Create}),
			init_fiber.NewApi("PUT", role_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Role, perm_constants.Update), roleApi.Update}),
			init_fiber.NewApi("DELETE", role_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.Role, perm_constants.Delete), roleApi.Delete}),

			init_fiber.NewApi("GET", user_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.User, perm_constants.Get), userApi.Get}),
			init_fiber.NewApi("GET", user_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.User, perm_constants.List), userApi.List}),
			init_fiber.NewApi("POST", user_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.User, perm_constants.Create), userApi.Create}),
			init_fiber.NewApi("PUT", user_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.User, perm_constants.Update), userApi.Update}),
			init_fiber.NewApi("DELETE", user_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.User, perm_constants.Delete), userApi.Delete}),

			init_fiber.NewApi("GET", user_group_role_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserGroupRole, perm_constants.Get), userGroupRoleApi.Get}),
			init_fiber.NewApi("GET", user_group_role_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserGroupRole, perm_constants.List), userGroupRoleApi.List}),
			init_fiber.NewApi("POST", user_group_role_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserGroupRole, perm_constants.Create), userGroupRoleApi.Create}),
			init_fiber.NewApi("PUT", user_group_role_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserGroupRole, perm_constants.Update), userGroupRoleApi.Update}),
			init_fiber.NewApi("DELETE", user_group_role_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserGroupRole, perm_constants.Delete), userGroupRoleApi.Delete}),

			init_fiber.NewApi("GET", user_product_url.Get, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserProduct, perm_constants.Get), userProductApi.Get}),
			init_fiber.NewApi("GET", user_product_url.List, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserProduct, perm_constants.List), userProductApi.List}),
			init_fiber.NewApi("POST", user_product_url.Create, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserProduct, perm_constants.Create), userProductApi.Create}),
			init_fiber.NewApi("PUT", user_product_url.Update, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserProduct, perm_constants.Update), userProductApi.Update}),
			init_fiber.NewApi("DELETE", user_product_url.Delete, []fiber.Handler{jwtValidator, token_util.AccessTokenValidator, token_util.HasResourcePerm(resource_constants.UserProduct, perm_constants.Delete), userProductApi.Delete}),
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
