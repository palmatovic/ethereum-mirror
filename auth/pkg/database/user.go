package database

import "time"

type User struct {
	UserId       int       `json:"user_id" gorm:"column:UserId;primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"column:Username;not null;uniqueKey:UsernameIdx"`
	Email        string    `json:"email" gorm:"column:Email;not null;uniqueKey:EmailIdx"`
	Password     string    `json:"-" gorm:"column:Password;not null"`
	Name         string    `json:"name" gorm:"column:Name;not null"`
	Surname      string    `json:"surname" gorm:"column:Surname;not null"`
	Dob          time.Time `json:"date_of_birth" gorm:"column:Dob;type:date;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
	SecretKey    string    `json:"-" gorm:"column:Secret"`
	TwoFAEnabled bool      `json:"-" gorm:"column:2FAEnabled;default:0"`
}

func (User) TableName() string {
	return "User"
}
