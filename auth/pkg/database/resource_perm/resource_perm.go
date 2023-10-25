package resource_perm

import (
	"auth/pkg/database/perm"
	"auth/pkg/database/resource"
	"time"
)

type ResourcePerm struct {
	ResourcePermId int64 `json:"resource_perm_id" gorm:"column:ResourcePermId;primaryKey;autoIncrement"`
	ResourceId     int64 `json:"resource_id" gorm:"column:ResourceId;uniqueIndex:ResourcePermIdx"`
	resource.Resource
	PermId int64 `json:"perm_id" gorm:"column:PermId;uniqueIndex:ResourcePermIdx"`
	perm.Perm
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
