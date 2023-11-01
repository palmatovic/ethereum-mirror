package user_product

import (
	"auth/pkg/database/product"
	"auth/pkg/database/user"
	"time"
)

type UserProduct struct {
	UserProductId        int64           `json:"user_product_id" gorm:"column:UserProductId;autoIncrement;primaryKey"`
	UserId               int64           `json:"user_id" gorm:"column:UserId;uniqueIndex:UserIdProductIdIdx"`
	User                 user.User       `json:"-"`
	ProductId            int64           `json:"product_id" gorm:"column:ProductId;uniqueIndex:UserIdProductIdIdx"`
	Product              product.Product `json:"-"`
	Enabled              bool            `json:"enabled" gorm:"column:UserEnabled;default:0"`
	ChangePassword       bool            `json:"change_password" gorm:"column:ChangePassword;default:0"`
	Password             string          `json:"password" gorm:"column:Password;not null"`
	PasswordExpirationAt time.Time       `json:"password_expiration_at" gorm:"column:PasswordExpirationAt"`
	PasswordExpired      bool            `json:"password_expired" gorm:"column:PasswordExpired;default:0"`
	MasterPasswordKey    string          `json:"master_password_key" gorm:"column:MasterKey"` // MasterPasswordKey used for reset lost/forgotten password
	TwoFAKey             string          `json:"two_fa_key" gorm:"column:2FAKey"`
	MasterTwoFAKey       string          `json:"master_two_fa_key" gorm:"column:MasterTwoFAKey"` // MasterTwoFAKey used for reset lost 2FA
	CreatedAt            time.Time       `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt            time.Time       `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserProduct) TableName() string {
	return "UserProduct"
}
