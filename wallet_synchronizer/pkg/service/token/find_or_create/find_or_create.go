package find_or_create

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	token_db "wallet-synchronizer/pkg/database/token"
	alchemy_token_metadata "wallet-synchronizer/pkg/service/alchemy/token_metadata"
	"wallet-synchronizer/pkg/service/goplus"
	token_create_service "wallet-synchronizer/pkg/service/token/create"
	token_get_service "wallet-synchronizer/pkg/service/token/get"
	token_update_service "wallet-synchronizer/pkg/service/token/update"
)

type Service struct {
	db              *gorm.DB
	contractAddress string
	alchemyApiKey   string
}

func NewService(db *gorm.DB, contractAddress string, alchemyApiKey string) *Service {
	return &Service{
		db:              db,
		contractAddress: contractAddress,
		alchemyApiKey:   alchemyApiKey,
	}
}

func (s *Service) FindOrCreateToken() (token *token_db.Token, err error) {
	var skipScam = false
	goplusService := goplus.NewService(s.contractAddress)
	if _, token, err = token_get_service.NewService(s.db, s.contractAddress).Get(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var newToken *token_db.Token
			if s.contractAddress == "ethereum" {
				newToken = &token_db.Token{
					TokenId:  s.contractAddress,
					Name:     "Ethereum",
					Symbol:   "ETH",
					Decimals: 18,
				}
			} else {
				var tokenMetadata alchemy_token_metadata.TokenMetadataResponse
				if tokenMetadata, err = alchemy_token_metadata.NewService(s.alchemyApiKey, s.contractAddress).TokenMetadata(); err != nil {
					return token, err
				}
				var logo string
				if len(tokenMetadata.Result.Logo) > 0 {
					logo, err = downloadLogo(tokenMetadata.Result.Logo)
					if err != nil {
						logrus.WithError(err).Errorf("cannot download logo from alchemy response for contract address %v", s.contractAddress)
					}
				}
				goplusResponse, errScam := goplusService.ScamCheck()
				if errScam != nil {
					return token, errScam
				}
				newToken = &token_db.Token{
					TokenId:        s.contractAddress,
					Name:           tokenMetadata.Result.Name,
					Symbol:         tokenMetadata.Result.Symbol,
					Decimals:       tokenMetadata.Result.Decimals,
					Logo:           logo,
					GoPlusResponse: goplusResponse,
				}
				skipScam = true
			}
			_, token, err = token_create_service.NewService(s.db, newToken).Create()
			if err != nil {
				return token, err
			}
		}
	}

	if !skipScam {
		goplusResponse, errScam := goplusService.ScamCheck()
		if errScam != nil {
			return token, errScam
		}
		token.GoPlusResponse = goplusResponse
		_, token, err = token_update_service.NewService(s.db, token).Update()
		if err != nil {
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
