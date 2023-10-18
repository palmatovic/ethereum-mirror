package perm

import "time"

type Perm struct {
	PermId    string    `json:"perm_id" gorm:"column:PermId"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
}

func (Perm) TableName() string {
	return "Perm"
}
