package address_status

import "transaction-extractor/pkg/database/address"

func (AddressStatus) TableName() string {
	return "AddressStatus"
}

type AddressStatus struct {
	AddressStatusId int    `gorm:"column:AddressStatusId;type:int;primaryKey;autoIncrement"`
	AddressId       string `gorm:"column:AddressId;varchar(1024)"`
	Address         address.Address
	Asset           string  `gorm:"column:Asset;varchar(1024);not null"`
	Symbol          string  `gorm:"column:Symbol;varchar(256);not null"`
	ContractAddress string  `gorm:"column:ContractAddress;varchar(1024);not null"`
	Quantity        float64 `gorm:"column:Quantity;not null"`
}
