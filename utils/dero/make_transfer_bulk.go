package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TransferInfo struct {
	Destination string `json:"destination"`
	Amount      int    `json:"amount"`
	SCID        string `json:"scid"`
}

func MakeBulkTransfer(transfers []TransferInfo) (string, error) {
	// Define payload data
	payloadData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "transfer",
		"params": map[string]interface{}{
			"ringsize":  8,
			"transfers": transfers,
			"scid":      DERO_SCID_STRING,
		},
	}

	// Marshal payload data to JSON
	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Define HTTP client and request
	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%s/json_rpc", DeroServerIP, deroTipsPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.SetBasicAuth(deroUser, deroPass)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Print response body
	fmt.Printf("Response Body: %s\n", responseBody)

	// Parse response JSON
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", fmt.Errorf("error decoding response JSON: %v", err)
	}

	if responseError, ok := response["error"].(map[string]interface{}); ok {
		return "", fmt.Errorf("DERO wallet error: %v", responseError)
	}

	if txID, ok := response["result"].(map[string]interface{})["txid"].(string); ok {
		return txID, nil
	}

	return "", fmt.Errorf("no transaction ID found in the response")
}
