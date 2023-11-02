package database

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"wallet-synchronizer/pkg/database/wallet"
	"wallet-synchronizer/pkg/service_util/aes"
)

type Service struct {
	aes256EncryptionKey *aes.Key
	dbFilepath          string
	tables              []interface{}
	ownWallet           string
}

func NewService(aes256EncryptionKey *aes.Key, dbFilepath string, ownWallet string, tables ...interface{}) *Service {
	return &Service{
		aes256EncryptionKey: aes256EncryptionKey,
		dbFilepath:          dbFilepath,
		tables:              tables,
		ownWallet:           ownWallet,
	}
}

func (s *Service) Init() (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlite.Open(s.dbFilepath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.AutoMigrate(
		s.tables...,
	)

	if err != nil {
		return nil, err
	}

	var ownWalletDb wallet.Wallet
	if err = db.Where("WalletId = ?", s.ownWallet).First(&ownWalletDb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ownWalletDb = wallet.Wallet{
				WalletId: s.ownWallet,
				Type:     false,
			}
			if err = db.Create(&ownWalletDb).Error; err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}
