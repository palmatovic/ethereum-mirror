package resource_perm

import (
	"auth/pkg/database/perm"
	"auth/pkg/database/resource"
	"time"
)

type ResourcePerm struct {
	ResourcePermId int64             `json:"resource_perm_id" gorm:"column:ResourcePermId;primaryKey;autoIncrement"`
	ResourceId     int64             `json:"resource_id" gorm:"column:ResourceId;uniqueIndex:ResourceIdPermIdIdx"`
	Resource       resource.Resource `json:"-"`
	PermId         string            `json:"perm_id" gorm:"column:PermId;uniqueIndex:ResourceIdPermIdIdx"`
	Perm           perm.Perm         `json:"-"`
	CreatedAt      time.Time         `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt      time.Time         `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (ResourcePerm) TableName() string {
	return "ResourcePerm"
}
