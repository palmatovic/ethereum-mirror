package wallet_ethereum_balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type WalletEthereumBalanceResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func GetWalletEthereumWalletBalance(wallet string, apiKey string) (walletEthereumBalanceResponse WalletEthereumBalanceResponse, err error) {
	// Alchemy URL
	var baseURL = fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", apiKey)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
		"params": []interface{}{
			wallet,
			"latest",
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
		return walletEthereumBalanceResponse, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		logrus.WithError(err).Error("received unexpected response status code from eth_getBalance request", response.StatusCode)
	}

	responseData, _ := io.ReadAll(response.Body)

	_ = json.Unmarshal(responseData, &walletEthereumBalanceResponse)
	return walletEthereumBalanceResponse, nil
}
