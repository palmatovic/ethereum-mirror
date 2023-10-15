package user

import "time"

type User struct {
	UserId               int       `json:"-" gorm:"column:UserId;primaryKey;autoIncrement"`
	Username             string    `json:"-" gorm:"column:Username;not null;uniqueKey:UsernameIdx"`
	UserEnabled          bool      `json:"-" gorm:"column:UserEnabled;default:0"`
	Email                string    `json:"-" gorm:"column:Email;not null;uniqueKey:EmailIdx"`
	Password             string    `json:"-" gorm:"column:Password;not null"`
	PasswordEnabled      bool      `json:"-" gorm:"column:PasswordEnabled;default:0"`
	PasswordExpirationAt time.Time `json:"-" gorm:"column:PasswordExpirationAt"`
	PasswordExpired      bool      `json:"-" gorm:"column:PasswordExpired;default:0"`
	PasswordAttempts     int       `json:"-" gorm:"column:PasswordAttempts"`
	MasterPasswordKey    string    `json:"-" gorm:"column:MasterKey"` // MasterPasswordKey used for reset lost/forgotten password
	Name                 string    `json:"-" gorm:"column:Name;not null"`
	Surname              string    `json:"-" gorm:"column:Surname;not null"`
	Dob                  time.Time `json:"-" gorm:"column:Dob;type:date;not null"`
	CreatedAt            time.Time `json:"-" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt            time.Time `json:"-" gorm:"column:UpdatedAt;autoUpdateTime"`
	TwoFAKey             string    `json:"-" gorm:"column:2FAKey"`
	TwoFAEnabled         bool      `json:"-" gorm:"column:2FAEnabled;default:0"`
	MasterTwoFAKey       string    `json:"-" gorm:"column:MasterTwoFAKey"` // MasterTwoFAKey used for reset lost 2FA
}

func (User) TableName() string {
	return "User"
}
