package database

import (
	company_db "auth/pkg/database/company"
	group_db "auth/pkg/database/group"
	group_role_db "auth/pkg/database/group_role"
	"auth/pkg/database/group_role_resource_perm"
	perm_db "auth/pkg/database/perm"
	product_db "auth/pkg/database/product"
	resource_db "auth/pkg/database/resource"
	resource_perm_db "auth/pkg/database/resource_perm"
	role_db "auth/pkg/database/role"
	user_db "auth/pkg/database/user"
	user_group_role_db "auth/pkg/database/user_group_role"
	user_product_db "auth/pkg/database/user_product"
	"auth/pkg/model/api/product/create"
	company_create_service "auth/pkg/service/company/create"
	group_create_service "auth/pkg/service/group/create"
	group_role_create_service "auth/pkg/service/group_role/create"
	group_role_resource_perm_create_service "auth/pkg/service/group_role_resource_perm/create"
	perm_create_service "auth/pkg/service/perm/create"
	product_create_service "auth/pkg/service/product/create"
	product_get_by_name_service "auth/pkg/service/product/get_by_name"
	resource_create_service "auth/pkg/service/resource/create"
	resource_perm_create_service "auth/pkg/service/resource_perm/create"
	role_create_service "auth/pkg/service/role/create"
	user_create_service "auth/pkg/service/user/create"
	user_group_role_create_service "auth/pkg/service/user_group_role/create"
	user_product_create_service "auth/pkg/service/user_product/create"
	"auth/pkg/service_util/aes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type Service struct {
	dbFilepath string
	tables     []interface{}
}

func NewService(dbFilepath string, tables ...interface{}) *Service {
	return &Service{dbFilepath: dbFilepath, tables: tables}
}

func (s *Service) Init() (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlite.Open(s.dbFilepath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.AutoMigrate(
		s.tables...,
	)

	if err != nil {
		return nil, err
	}

	var product *product_db.Product
	if _, product, err = product_get_by_name_service.NewService(tx, "auth").Get(); err == nil {
		return db, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, product, err = product_create_service.NewService(tx, &create.Product{
		Name:        "auth",
		Description: "auth product",
		SslSetup: create.SslSetup{
			Company:     "auth-company",
			Province:    "Rome",
			Country:     "IT",
			CompanyUnit: "ENG",
			CommonName:  "*.auth-company.com",
			Locality:    "Rome",
			AltDNS:      "*.auth-company.com",
		},
		JwtConfig: create.JwtConfig{
			AccessToken:  create.Token{ExpiresInMinutes: int64(4 * 60)},
			RefreshToken: create.Token{ExpiresInMinutes: int64(8 * 60)},
		},
	}).Create(); err != nil {
		return nil, err
	}

	var company *company_db.Company
	if _, company, err = company_create_service.NewService(tx, &company_db.Company{
		Name: "auth-company",
	}).Create(); err != nil {
		return nil, err
	}

	var group *group_db.Group
	if _, group, err = group_create_service.NewService(tx, &group_db.Group{
		Name:      "admin-group",
		ProductId: product.ProductId,
		CompanyId: company.CompanyId,
	}).Create(); err != nil {
		return nil, err
	}

	var role *role_db.Role
	if _, role, err = role_create_service.NewService(tx, &role_db.Role{
		Name: "admin-role",
	}).Create(); err != nil {
		return nil, err
	}

	var groupRole *group_role_db.GroupRole
	if _, groupRole, err = group_role_create_service.NewService(tx, &group_role_db.GroupRole{
		GroupId: group.GroupId,
		RoleId:  role.RoleId,
	}).Create(); err != nil {
		return nil, err
	}

	var resources = new([]resource_db.Resource)
	for _, value := range []string{
		"login",
		"logout",
		"otp",
		"product",
		"company",
		"group",
		"role",
		"group_role",
		"resource",
		"perm",
		"resource_perm",
		"group_role_resource_perm",
		"user_product",
		"user_resource_perm",
		"change_password",
	} {
		var resource *resource_db.Resource
		if _, resource, err = resource_create_service.NewService(tx, &resource_db.Resource{
			Name: value,
		}).Create(); err != nil {
			return nil, err
		}
		*resources = append(*resources, *resource)
	}

	var perms = new([]perm_db.Perm)
	for _, value := range []string{
		"create",
		"update",
		"get",
		"list",
		"delete",
	} {
		var perm *perm_db.Perm
		if _, perm, err = perm_create_service.NewService(tx, &perm_db.Perm{
			PermId: value,
		}).Create(); err != nil {
			return nil, err
		}
		*perms = append(*perms, *perm)
	}

	var resourcePerms = new([]resource_perm_db.ResourcePerm)
	for _, res := range *resources {
		for _, per := range *perms {
			var resourcePerm *resource_perm_db.ResourcePerm
			if _, resourcePerm, err = resource_perm_create_service.NewService(tx, &resource_perm_db.ResourcePerm{
				ResourceId: res.ResourceId,
				PermId:     per.PermId,
			}).Create(); err != nil {
				return nil, err
			}
			*resourcePerms = append(*resourcePerms, *resourcePerm)
		}
	}

	var groupRoleResourcePerms = new([]group_role_resource_perm.GroupRoleResourcePerm)
	for _, value := range *resourcePerms {
		var groupRoleResourcePerm *group_role_resource_perm.GroupRoleResourcePerm
		if _, groupRoleResourcePerm, err = group_role_resource_perm_create_service.NewService(tx, &group_role_resource_perm.GroupRoleResourcePerm{
			GroupRoleId:    groupRole.GroupRoleId,
			ResourcePermId: value.ResourcePermId,
		}).Create(); err != nil {
			return nil, err
		}
		*groupRoleResourcePerms = append(*groupRoleResourcePerms, *groupRoleResourcePerm)
	}

	var user *user_db.User
	if _, user, err = user_create_service.NewService(tx, &user_db.User{
		CompanyId:   company.CompanyId,
		Username:    "auth",
		Name:        "auth",
		Surname:     "auth",
		DateOfBirth: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()),
	}).Create(); err != nil {
		return nil, err
	}

	if _, _, err = user_group_role_create_service.NewService(tx, &user_group_role_db.UserGroupRole{
		UserId:      user.UserId,
		GroupRoleId: groupRole.GroupRoleId,
	}).Create(); err != nil {
		return nil, err
	}

	// generate two fa key
	masterPwdKey, err := aes.NewService().NewAES256Key()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%x\n", *masterPwdKey)
	masterTwoFAKey, err := aes.NewService().NewAES256Key()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%x\n", *masterTwoFAKey)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "auth",
		AccountName: "auth",
		Algorithm:   otp.AlgorithmSHA256,
	})
	if err != nil {
		return nil, err
	}

	file, err := os.Create("admin_product-auth.json")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	// Crea un encoder JSON
	encoder := json.NewEncoder(file)

	userProduct := &user_product_db.UserProduct{
		UserProductId:        0,
		UserId:               user.UserId,
		User:                 user_db.User{},
		ProductId:            product.ProductId,
		Product:              product_db.Product{},
		Enabled:              true,
		Password:             "admin-password",
		PasswordExpirationAt: time.Now().Add(time.Hour * 24 * 365),
		PasswordExpired:      false,
		MasterPasswordKey:    *masterPwdKey,
		TwoFAKey:             key.Secret(),
		MasterTwoFAKey:       *masterTwoFAKey,
	}

	if err = encoder.Encode(userProduct); err != nil {
		return nil, err
	}
	sha256MasterPwdKey := sha256.Sum256(userProduct.MasterPasswordKey)
	userProduct.MasterPasswordKey = sha256MasterPwdKey[:]
	fmt.Printf("%x\n", userProduct.MasterPasswordKey)

	sha256MasterTwoFAKey := sha256.Sum256(userProduct.MasterTwoFAKey)
	userProduct.MasterTwoFAKey = sha256MasterTwoFAKey[:]
	fmt.Printf("%x\n", userProduct.MasterTwoFAKey)

	if _, _, err = user_product_create_service.NewService(tx, userProduct).Create(); err != nil {
		return nil, err
	}

	return db, nil
}
