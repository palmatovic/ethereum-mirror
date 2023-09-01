package address_status

import "transaction-extractor/pkg/database/address"

func (AddressStatus) TableName() string {
	return "AddressStatus"
}

type AddressStatus struct {
	AddressId            string `gorm:"column:AddressId;primaryKey;varchar(1024)"`
	Address              address.Address
	TokenName            string  `gorm:"column:TokenName;varchar(1024);not null"`
	TokenSymbol          string  `gorm:"column:TokenSymbol;varchar(256);not null"`
	TokenContractAddress string  `gorm:"column:TokenContractAddress;primaryKey;varchar(1024);not null"`
	TokenAmount          float64 `gorm:"column:TokenAmount;not null"`
}
