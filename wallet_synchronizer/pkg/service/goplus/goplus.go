package goplus

import (
	"errors"
	"github.com/GoPlusSecurity/goplus-sdk-go/api/token"
	"github.com/GoPlusSecurity/goplus-sdk-go/pkg/errorcode"
	"github.com/GoPlusSecurity/goplus-sdk-go/pkg/gen/models"
)

type Service struct {
	tokenAddress string
}

func NewService(tokenAddress string) *Service {
	return &Service{tokenAddress: tokenAddress}
}

func (s *Service) ScamCheck() (isScam bool, err error) {
	tokenSecurity := token.NewTokenSecurity(nil)
	chainId := "1"
	contractAddresses := []string{s.tokenAddress}
	data, err := tokenSecurity.Run(chainId, contractAddresses)
	if err != nil {
		return false, err
	}
	if data.Payload.Code != errorcode.SUCCESS {
		return false, err
	}
	response, ok := data.Payload.Result[s.tokenAddress]

	if ok {
		isScam = doScamCheck(&response)
		return isScam, nil
	} else {
		return false, errors.New("result not contains token address")
	}
}

func doScamCheck(response *models.ResponseWrapperTokenSecurityResultAnon) bool {

	if response.TrustList == "1" {
		return false
	}

	if response.IsTrueToken == "0" {
		return true
	}

	riskyItems, _ := riskWarningCount(response)

	return riskyItems >= 1
}

func riskWarningCount(response *models.ResponseWrapperTokenSecurityResultAnon) (risky int, attention int) {

	if response.IsOpenSource == "0" {
		attention++
	}

	if response.CanTakeBackOwnership == "1" {
		attention++
	}

	if response.IsMintable == "1" {
		attention++
	}

	if response.OwnerChangeBalance == "1" {
		risky++
	}

	if response.HiddenOwner == "1" {
		attention++
	}

	if response.CannotSellAll == "1" {
		risky++
	}

	if response.TransferPausable == "1" {
		attention++
	}

	// honeypot
	if response.IsHoneypot == "1" || response.HoneypotWithSameCreator == "1" {
		risky++
	}

	if response.TradingCooldown == "1" {
		attention++
	}

	if response.IsWhitelisted == "1" {
		attention++
	}

	if response.AntiWhaleModifiable == "1" {
		attention++
	}

	if response.IsAntiWhale == "1" {
		attention++
	}

	if response.IsBlacklisted == "1" {
		attention++
	}

	if response.SlippageModifiable == "1" {
		attention++
	}

	if response.TransferPausable == "1" {
		attention++
	}

	return risky, attention

}
