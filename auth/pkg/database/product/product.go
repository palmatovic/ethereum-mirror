package product

import "time"

type Product struct {
	ProductId string    `json:"product_id" gorm:"column:ProductId;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Product) TableName() string {
	return "Product"
}
