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

func (s *Service) ScamCheck() (response models.ResponseWrapperTokenSecurityResultAnon, err error) {
	tokenSecurity := token.NewTokenSecurity(nil)
	chainId := "1"
	contractAddresses := []string{s.tokenAddress}
	data, err := tokenSecurity.Run(chainId, contractAddresses)
	if err != nil {
		return response, err
	}
	if data.Payload.Code != errorcode.SUCCESS {
		return response, err
	}
	var ok bool
	response, ok = data.Payload.Result[s.tokenAddress]
	doScamCheck(&response)
	if ok {
		return response, nil
	} else {
		return response, errors.New("result not contains token address")
	}
}

func doScamCheck(response *models.ResponseWrapperTokenSecurityResultAnon) bool {
	riskyItems, _ := riskWarningCount(response)

	return riskyItems >= 1
}

func riskWarningCount(response *models.ResponseWrapperTokenSecurityResultAnon) (risky int, attention int) {
	if response.IsMintable == "1" {
		attention++
	}

	if response.OwnerChangeBalance == "1" {
		risky++
	}

	if response.HiddenOwner == "1" {
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
