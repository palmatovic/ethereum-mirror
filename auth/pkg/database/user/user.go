package user

import "time"

type User struct {
	UserId    int       `json:"-" gorm:"column:UserId;primaryKey;autoIncrement"`
	Username  string    `json:"-" gorm:"column:Username;not null;uniqueKey:UsernameIdx"`
	Email     string    `json:"-" gorm:"column:Email;not null;uniqueKey:EmailIdx"`
	Name      string    `json:"-" gorm:"column:Name;not null"`
	Surname   string    `json:"-" gorm:"column:Surname;not null"`
	Dob       time.Time `json:"-" gorm:"column:Dob;type:date;not null"`
	CreatedAt time.Time `json:"-" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (User) TableName() string {
	return "User"
}
