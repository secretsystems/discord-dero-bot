package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func MakeTransfer(address string) {
	// Define payload data
	addr := address
	amnt := 2
	scid := "0000000000000000000000000000000000000000000000000000000000000000"
	payloadData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "transfer",
		"params": map[string]interface{}{
			"transfers": []map[string]interface{}{
				{
					"scid":        scid,
					"destination": addr,
					"amount":      amnt,
				},
			},
			"ringsize": 16,
		},
	}

	// Marshal payload data to JSON
	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Define HTTP client and request
	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%s/json_rpc", os.Getenv("DERO_SERVER_IP"), os.Getenv("DERO_WALLET_PORT"))
	username := os.Getenv("USER")
	password := os.Getenv("PASS")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print response body
	fmt.Printf("Response Body: %s\n", responseBody)

	// Parse response JSON
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Fatalf("Error decoding response JSON: %v", err)
	}

	// Assuming you want to print the parsed response
	fmt.Printf("Parsed Response: %+v\n", response)
}
