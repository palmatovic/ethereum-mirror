package product

import "time"

type Product struct {
	ProductId                         int64     `json:"product_id" gorm:"column:ProductId;primaryKey;autoIncrement"`
	Name                              string    `json:"name" gorm:"column:Name"`
	Description                       string    `json:"description" gorm:"column:Description"`
	ServerCert                        []byte    `json:"server_cert" gorm:"column:ServerCert"`
	ServerKey                         []byte    `json:"server_key" gorm:"column:ServerKey"`
	CaCert                            []byte    `json:"ca_cert" gorm:"column:CaCert"`
	CaKey                             []byte    `json:"ca_key" gorm:"column:CaKey"`
	ExpiredSsl                        bool      `json:"expired_ssl" gorm:"column:ExpiredSsl;default=0"`
	RSAPrivateKey                     []byte    `json:"rsa_private_key" gorm:"column:RSAPrivateKey"`
	RSAPublicKey                      []byte    `json:"rsa_public_key" gorm:"column:RSAPublicKey"`
	JwtConfig                         []byte    `json:"jwt_config" gorm:"column:JWTConfig"`
	JwtExpirationDate                 time.Time `json:"jwt_expiration_date" gorm:"column:JwtExpirationDate"`
	AES256EncryptionKey               []byte    `json:"aes_256_encryption_key" gorm:"column:AES256EncryptionKey"`
	AES256EncryptionKeyExpirationDate time.Time `json:"aes_256_encryption_key_expiration_date" gorm:"column:AES256EncryptionKeyExpirationDate"`
	CreatedAt                         time.Time `json:"created_at" gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt                         time.Time `json:"updated_at" gorm:"column:UpdatedAt;autoUpdateTime"`
}
