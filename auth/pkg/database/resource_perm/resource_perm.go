package resource_perm

import (
	"auth/pkg/database/perm"
	"auth/pkg/database/resource"
	"time"
)

type ResourcePerm struct {
	ResourcePermId int `json:"resource_perm_id" gorm:"ResourcePermId;autoIncrement;primaryKey"`
	ResourceId     int `json:"resource_id" gorm:"column:ResourceId;uniqueIndex:ResourcePermIdx"`
	resource.Resource
	PermId string `json:"perm_id" gorm:"PermId;uniqueIndex:ResourcePermIdx"`
	perm.Perm
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (ResourcePerm) TableName() string {
	return "ResourcePerm"
}
