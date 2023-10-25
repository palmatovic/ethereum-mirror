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
	Enabled              bool      `json:"enabled" gorm:"column:UserEnabled;default:0"`
	Password             string    `json:"password" gorm:"column:Password;not null"`
	PasswordExpirationAt time.Time `json:"password_expiration_at" gorm:"column:PasswordExpirationAt"`
	PasswordExpired      bool      `json:"password_expired" gorm:"column:PasswordExpired;default:0"`
	MasterPasswordKey    string    `json:"master_password_key" gorm:"column:MasterKey"` // MasterPasswordKey used for reset lost/forgotten password
	TwoFAKey             string    `json:"two_fa_key" gorm:"column:2FAKey"`
	MasterTwoFAKey       string    `json:"master_two_fa_key" gorm:"column:MasterTwoFAKey"` // MasterTwoFAKey used for reset lost 2FA
	CreatedAt            time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
