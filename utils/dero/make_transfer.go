package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func MakeTransfer(address string, amnt int, comment string) {
	// Define payload data
	addr := address

	scid := "0000000000000000000000000000000000000000000000000000000000000000"
	payloadData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "transfer",
		"params": map[string]interface{}{
			"ringsize": 16,
			"transfers": []map[string]interface{}{
				{
					"scid":        scid,
					"destination": addr,
					"amount":      amnt,
					"payload": map[string]interface{}{
						"name":     "C",
						"datatype": "S",
						"value":    comment,
					},
				},
			},
		},
	}

	// Marshal payload data to JSON
	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Define HTTP client and request
	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%s/json_rpc", deroServerIP, deroWalletPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.SetBasicAuth(deroUser, deroPass)
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
