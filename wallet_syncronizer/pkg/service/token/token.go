package token

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	token_db "wallet-syncronizer/pkg/database/token"
	alchemy_token_metadata "wallet-syncronizer/pkg/service/alchemy/token_metadata"
	"wallet-syncronizer/pkg/service/goplus"
)

func FindOrCreateToken(db *gorm.DB, contractAddress string, alchemyApiKey string) (token token_db.Token, err error) {
	var skipScam = false
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
			goplusResponse, errScam := goplus.ScamCheck(contractAddress)
			if errScam != nil {
				return token, errScam
			}
			token = token_db.Token{
				TokenId:  contractAddress,
				Name:     tokenMetadata.Result.Name,
				Symbol:   tokenMetadata.Result.Symbol,
				Decimals: tokenMetadata.Result.Decimals,
				Logo:     logo,
				GoPlusResponse: func() []byte {
					b, _ := json.Marshal(goplusResponse)
					return b
				}(),
			}
			if err = db.Create(&token).Error; err != nil {
				return token, err
			}
			skipScam = true
		}
	}

	if !skipScam {
		goplusResponse, errScam := goplus.ScamCheck(contractAddress)
		if errScam != nil {
			return token, errScam
		}
		token.GoPlusResponse = func() []byte {
			b, _ := json.Marshal(goplusResponse)
			return b
		}()
		if err = db.Where("TokenId = ?", token.TokenId).Updates(&token).Error; err != nil {
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

func GetToken(db *gorm.DB, tokenId string) (status int, token *token_db.Token, err error) {
	token = new(token_db.Token)
	if err = db.Where("TokenId = ?", tokenId).First(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.StatusNotFound, nil, err
		} else {
			return fiber.StatusInternalServerError, nil, err
		}
	}
	return fiber.StatusOK, token, nil
}
