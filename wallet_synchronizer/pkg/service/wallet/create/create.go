package create

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
	wallet_db "wallet-synchronizer/pkg/database/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	db     *gorm.DB
	wallet *wallet_db.Wallet
}

func NewService(db *gorm.DB, wallet *wallet_db.Wallet) *Service {
	return &Service{
		db:     db,
		wallet: wallet,
	}
}

func (s *Service) Create() (status int, wallet *wallet_db.Wallet, err error) {
	// check if wallet is correct or exists on blockchain before creating
	var ok bool
	if ok, err = isValidWallet(s.wallet.WalletId); err != nil || !ok {
		if err != nil {
			return fiber.StatusInternalServerError, nil, err
		}
		return fiber.StatusBadRequest, nil, errors.New("invalid address type, must be a wallet address")
	}
	if err = s.db.Create(s.wallet).Error; err != nil {
		return fiber.StatusInternalServerError, nil, err
	}
	return fiber.StatusOK, s.wallet, nil
}

func isValidWallet(address string) (bool, error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if !re.MatchString(address) {
		return false, errors.New("invalid wallet address: doesn't match address regex")
	}

	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		return false, err
	}

	// a random user account address
	checkAddress := common.HexToAddress(address)
	bytecode, err := client.CodeAt(context.Background(), checkAddress, nil)
	if err != nil {
		return false, err
	}

	return len(bytecode) <= 0, nil
}
