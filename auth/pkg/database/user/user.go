package user

import (
	"auth/pkg/database/company"
	"time"
)

type User struct {
	UserId    int64 `json:"user_id" gorm:"column:UserId;primaryKey;autoIncrement"`
	CompanyId int64 `json:"company_id" gorm:"column:CompanyId"`
	company.Company
	Username    string    `json:"username" gorm:"column:Username;uniqueIndex:UserUsernameIdx"`
	Name        string    `json:"name" gorm:"column:Name"`
	Surname     string    `json:"surname" gorm:"column:Surname"`
	DateOfBirth string    `json:"date_of_birth" gorm:"column:DateOfBirth;type:date"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
