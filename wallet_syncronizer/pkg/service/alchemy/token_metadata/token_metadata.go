package token_metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type TokenMetadataResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  Result `json:"result"`
}
type Result struct {
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
}

type Service struct {
	apiKey          string
	contractAddress string
}

func NewService(apiKey string, contractAddress string) *Service {
	return &Service{
		apiKey:          apiKey,
		contractAddress: contractAddress,
	}
}

func (s *Service) TokenMetadata() (tokenMetadataResponse TokenMetadataResponse, err error) {
	// Alchemy URL
	var baseURL = fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", s.apiKey)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenMetadata",
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
		"params": []interface{}{
			s.contractAddress,
		},
		"id": 1,
	}

	payload, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.WithError(err).Error("cannot connect to server")
		return tokenMetadataResponse, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		logrus.WithError(err).Error("received unexpected response status code from alchemy_getTokenMetadata request", response.StatusCode)
	}

	responseData, _ := io.ReadAll(response.Body)

	_ = json.Unmarshal(responseData, &tokenMetadataResponse)
	return tokenMetadataResponse, nil
}
