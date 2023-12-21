package monero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gabstv/httpdigest"
)

type BalanceResponse struct {
	Result struct {
		Balance uint64 `json:"balance"`
	} `json:"result"`
}

func GetWalletBalance() (uint64, error) {
	log.Println("Fetching wallet balance...")

	accountIndex := uint(0)
	addressIndices := []uint{0, 1}
	allAccounts := false
	strict := false

	params := map[string]interface{}{
		"account_index":   accountIndex,
		"address_indices": addressIndices,
		"all_accounts":    allAccounts,
		"strict":          strict,
	}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_balance",
		"params":  params,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return 0, fmt.Errorf("error marshaling JSON data: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/json_rpc", moneroServerIP, moneroWalletPort)
	log.Printf("URL: %s", url)

	// Create an HTTP client with Digest Authentication
	client := http.Client{
		Transport: &httpdigest.Transport{
			Username:  moneroUser,
			Password:  moneroPass,
			Transport: http.DefaultTransport,
		},
	}

	// Create the HTTP request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP POST request: %v", err)
		return 0, fmt.Errorf("error sending HTTP POST request: %v", err)
	}
	defer response.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var balanceResponse BalanceResponse
	err = json.Unmarshal(responseBody, &balanceResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return 0, fmt.Errorf("error decoding response JSON: %v", err)
	}

	log.Println("Wallet balance fetched successfully.")

	return balanceResponse.Result.Balance, nil
}
