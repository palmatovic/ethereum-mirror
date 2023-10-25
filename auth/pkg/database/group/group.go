package group

import (
	"auth/pkg/database/company"
	"auth/pkg/database/product"
	"time"
)

type Group struct {
	GroupId   int64  `json:"group_id" gorm:"column:GroupId;primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"column:Name"`
	ProductId int64  `json:"product_id" gorm:"column:ProductId;uniqueIndex:ProductCompanyIdx"`
	product.Product
	CompanyId int64 `json:"company_id" gorm:"column:CompanyId;uniqueIndex:ProductCompanyIdx"`
	company.Company
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
