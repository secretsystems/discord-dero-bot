package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func MakeTransfer() {
	addr := "destination_address" // Replace with actual destination address
	// dstPort := "destination_port" // Replace with actual destination port
	amnt := 2

	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "transfer",
		"params": map[string]interface{}{
			"scid":        "00000000000000000000000000000000",
			"destination": addr,
			"amount":      amnt,
			"payload_rpc": []map[string]interface{}{
				{
					"name":     "C",
					"datatype": "S",
					"value":    "Secret Message: luv u",
				},
				{
					"name":     "D",
					"datatype": "U",
					"value":    7331,
				},
				{
					"name":     "S",
					"datatype": "U",
					"value":    1337,
				},
			},
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/json_rpc", os.Getenv("DERO_SERVER_IP"), os.Getenv("DERO_WALLET_PORT"))
	username := os.Getenv("USER")
	password := os.Getenv("PASS")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{} // Use a new instance of http.Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	// Assuming you want to print the response
	fmt.Printf("Response: %+v\n", response)
}
