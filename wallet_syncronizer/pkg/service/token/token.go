package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	token_db "wallet-syncronizer/pkg/database/token"
	alchemy_token_metadata "wallet-syncronizer/pkg/service/alchemy/token_metadata"
	"wallet-syncronizer/pkg/service/goplus"
)

func FindOrCreateToken(db *gorm.DB, contractAddress string, alchemyApiKey string, browser playwright.Browser) (token token_db.Token, err error) {
	if err = db.Where("TokenId = ?", contractAddress).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var tokenMetadata alchemy_token_metadata.TokenMetadataResponse
			tokenMetadata, err = alchemy_token_metadata.GetTokenMetadata(alchemyApiKey, contractAddress)
			if err != nil {
				return token, err
			}
			var logo string
			if len(tokenMetadata.Result.Logo) > 0 {
				logo, err = downloadLogo(tokenMetadata.Result.Logo)
				if err != nil {
					logrus.WithError(err).Errorf("cannot download logo from alchemy response for contract address %v", contractAddress)
				}
			}
			token = token_db.Token{
				TokenId:  contractAddress,
				Name:     tokenMetadata.Result.Name,
				Symbol:   tokenMetadata.Result.Symbol,
				Decimals: tokenMetadata.Result.Decimals,
				Logo:     logo,
			}
			if err = db.Create(&token).Error; err != nil {
				return token, err
			}
		}
	}

	risk, warn, err := goplus.IsScam(token.TokenId, browser)
	if err != nil {
		return token, err
	}

	if token.RiskScam != risk || token.WarningScam != warn {
		if token.RiskScam != risk {
			token.RiskScam = risk
		}
		if token.WarningScam != warn {
			token.WarningScam = warn
		}
		if err = db.Updates(&token).Error; err != nil {
			return token, err
		}
	}

	return token, nil
}

func downloadLogo(logoUrl string) (string, error) {
	response, errGet := http.Get(logoUrl)
	if errGet != nil {
		return "", errGet
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status code: %d", response.StatusCode)
	}
	imageBytes, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		return "", errRead
	}
	return base64.StdEncoding.EncodeToString(imageBytes), nil
}
