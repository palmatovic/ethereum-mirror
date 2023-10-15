package resource

import (
	"auth/pkg/database/product"
	"time"
)

type Resource struct {
	ResourceId string `json:"resource_id" gorm:"column:ResourceId;autoIncrement;primaryKey"`
	Name       string `json:"name" gorm:"column:Name;not mull"`
	ProductId  string `json:"product_id" gorm:"column:ProductId"`
	product.Product
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Resource) TableName() string {
	return "Resource"
}
