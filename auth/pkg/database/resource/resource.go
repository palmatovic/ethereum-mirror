package resource

import "time"

type Resource struct {
	ResourceId int64     `json:"resource_id" gorm:"column:ResourceId;primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"column:Name"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Resource) TableName() string {
	return "Resource"
}
