package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Rest of your FetchDeroTransfers function and other code ...

func FetchDeroTransfers() ([]byte, error) {
	serverIP := os.Getenv("DERO_SERVER_IP")
	serverPort := os.Getenv("DERO_WALLET_PORT")
	url := fmt.Sprintf("http://%s:%s/json_rpc", serverIP, serverPort)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "GetTransfers",
		"params": map[string]interface{}{
			"coinbase": true,
			"out":      true,
			"in":       true,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(os.Getenv("USER"), os.Getenv("PASS")) // Use values from .env
	req.Header.Set("Content-type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
