package role

import (
	"time"
)

type Role struct {
	RoleId    int64     `json:"role_id" gorm:"column:RoleId;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:Name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
