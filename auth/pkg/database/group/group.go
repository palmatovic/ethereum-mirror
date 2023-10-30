package group

import (
	"auth/pkg/database/company"
	"auth/pkg/database/product"
	"time"
)

type Group struct {
	GroupId   int64           `json:"group_id" gorm:"column:GroupId;primaryKey;autoIncrement"`
	Name      string          `json:"name" gorm:"column:Name"`
	ProductId int64           `json:"product_id" gorm:"column:ProductId;uniqueIndex:ProductIdCompanyIdIdx"`
	Product   product.Product `json:"-"`
	CompanyId int64           `json:"company_id" gorm:"column:CompanyId;uniqueIndex:ProductIdCompanyIdIdx"`
	Company   company.Company `json:"-"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Group) TableName() string {
	return "Group"
}
