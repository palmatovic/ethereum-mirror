package user_product

import (
	"auth/pkg/database/product"
	"auth/pkg/database/user"
	"time"
)

type UserProduct struct {
	UserProductId int `json:"user_product_id" gorm:"UserProductId;autoIncrement;primaryKey"`
	UserId        int `json:"user_id" gorm:"column:UserId;uniqueIndex:UserProductIdx"`
	user.User
	ProductId string `json:"product_id" gorm:"ProductId;uniqueIndex:UserProductIdx"`
	product.Product
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserProduct) TableName() string {
	return "UserProduct"
}
