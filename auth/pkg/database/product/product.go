package product

import "time"

// Auth non ha bisogno di salvare Rsa256PK
type Product struct {
	ProductId                   string    `json:"product_id" gorm:"column:ProductId;primaryKey"`
	Rsa256PK                    string    `json:"-" gorm:"column:Rsa256PK"`
	AccessTokenDurationMinutes  int       `json:"access_token_duration_minutes" gorm:"column:AccessTokenDurationMinutes"`
	RefreshTokenDurationMinutes int       `json:"refresh_token_duration_minutes" gorm:"column:RefreshTokenDurationMinutes"`
	CreatedAt                   time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt                   time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Product) TableName() string {
	return "Product"
}
