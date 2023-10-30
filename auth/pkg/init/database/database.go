package database

import (
	company_db "auth/pkg/database/company"
	group_db "auth/pkg/database/group"
	group_role_db "auth/pkg/database/group_role"
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
	product_create_service "auth/pkg/service/product/create"
	product_get_by_name_service "auth/pkg/service/product/get_by_name"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	err = db.AutoMigrate(
		s.tables...,
	)

	var product *product_db.Product
	if _, product, err = product_get_by_name_service.NewService(db, "auth").Get(); err == nil {
		return db, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tx := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if _, product, err = product_create_service.NewService(db, &create.Product{
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
			AccessToken:  create.Token{ExpiresInMinutes: int64(time.Minute * 4 * 60)},
			RefreshToken: create.Token{ExpiresInMinutes: int64(time.Minute * 8 * 60)},
		},
	}).Create(); err != nil {
		return nil, err
	}

	var company *company_db.Company
	if _, company, err = company_create_service.NewService(db, &company_db.Company{
		Name: "auth-company",
	}).Create(); err != nil {
		return nil, err
	}

	var group *group_db.Group
	if _, group, err = group_create_service.NewService(db, &group_db.Group{}).Create(); err != nil {
		return nil, err
	}

	var role *role_db.Role
	if _, role, err = role_create_service.NewService(db, &role_db.Role{}).Create(); err != nil {
		return nil, err
	}

	var groupRole *group_role_db.GroupRole
	if _, groupRole, err = group_role_create_service.NewService(db, &group_role_db.GroupRole{}).Create(); err != nil {
		return nil, err
	}

	var resource *resource_db.Resource
	if _, resource, err = resource_create_service.NewService(db, &resource_db.Resource{}).Create(); err != nil {
		return nil, err
	}

	var perm *perm_db.Perm
	if _, perm, err = perm_create_service.NewService(db, &perm_db.Perm{}).Create(); err != nil {
		return nil, err
	}

	var resourcePerm *resource_perm_db.ResourcePerm
	if _, resourcePerm, err = resource_perm_create_service.NewService(db, &resource_perm_db.ResourcePerm{}).Create(); err != nil {
		return nil, err
	}

	var user *user_db.User
	if _, user, err = user_create_service.NewService(db, &user_db.User{}).Create(); err != nil {
		return nil, err
	}

	var userGroupRole *user_group_role_db.UserGroupRole
	if _, userGroupRole, err = user_group_role_create_service.NewService(db, &user_group_role_db.UserGroupRole{}).Create(); err != nil {
		return nil, err
	}

	var userProduct *user_product_db.UserProduct
	if _, userProduct, err = user_product_create_service.NewService(db, &user_product_db.UserProduct{}).Create(); err != nil {
		return nil, err
	}

	// creare gruppi, ruoli, risorse ecc

	return db, nil
}
