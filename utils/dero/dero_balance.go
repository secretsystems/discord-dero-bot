package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeroBalanceResponse struct {
	Result struct {
		Balance         uint64 `json:"balance"`
		UnlockedBalance uint64 `json:"unlocked_balance"`
	} `json:"result"`
}

func GetDeroWalletBalance() (uint64, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "GetBalance",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("error marshaling JSON data: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/json_rpc", deroServerIP, deroWalletPort)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}

	request.SetBasicAuth(deroUser, deroPass) // Use values from .env
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP POST request: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var balanceResponse DeroBalanceResponse
	err = json.Unmarshal(responseBody, &balanceResponse)
	if err != nil {
		return 0, fmt.Errorf("error decoding response JSON: %v", err)
	}

	return balanceResponse.Result.Balance, nil
}
