package perm

import "time"

type Perm struct {
	PermId    string    `json:"perm_id" gorm:"PermId;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"UpdatedAt;autoUpdateTime"`
}
