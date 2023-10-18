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
	UserEnabled          bool      `json:"-" gorm:"column:UserEnabled;default:0"`
	Password             string    `json:"-" gorm:"column:Password;not null"`
	PasswordEnabled      bool      `json:"-" gorm:"column:PasswordEnabled;default:0"`
	PasswordExpirationAt time.Time `json:"-" gorm:"column:PasswordExpirationAt"`
	PasswordExpired      bool      `json:"-" gorm:"column:PasswordExpired;default:0"`
	PasswordAttempts     int       `json:"-" gorm:"column:PasswordAttempts"`
	MasterPasswordKey    string    `json:"-" gorm:"column:MasterKey"` // MasterPasswordKey used for reset lost/forgotten password
	TwoFAKey             string    `json:"-" gorm:"column:2FAKey"`
	TwoFAEnabled         bool      `json:"-" gorm:"column:2FAEnabled;default:0"`
	MasterTwoFAKey       string    `json:"-" gorm:"column:MasterTwoFAKey"` // MasterTwoFAKey used for reset lost 2FA
	CreatedAt            time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserProduct) TableName() string {
	return "UserProduct"
}
