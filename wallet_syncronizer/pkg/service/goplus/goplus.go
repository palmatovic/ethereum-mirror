package goplus

import (
	"errors"
	"github.com/GoPlusSecurity/goplus-sdk-go/api/token"
	"github.com/GoPlusSecurity/goplus-sdk-go/pkg/errorcode"
	"github.com/GoPlusSecurity/goplus-sdk-go/pkg/gen/models"
)

func ScamCheck(tokenAddress string) (response models.ResponseWrapperTokenSecurityResultAnon, err error) {
	tokenSecurity := token.NewTokenSecurity(nil)
	chainId := "1"
	contractAddresses := []string{tokenAddress}
	data, err := tokenSecurity.Run(chainId, contractAddresses)
	if err != nil {
		return response, err
	}
	if data.Payload.Code != errorcode.SUCCESS {
		return response, err
	}
	var ok bool
	response, ok = data.Payload.Result[tokenAddress]

	if ok {
		return response, nil
	} else {
		return response, errors.New("result not contains token address")
	}
}
