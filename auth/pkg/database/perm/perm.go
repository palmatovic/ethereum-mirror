package perm

import "time"

type Perm struct {
	PermId    string    `json:"perm_id" gorm:"column:PermId;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Perm) TableName() string {
	return "Perm"
}
