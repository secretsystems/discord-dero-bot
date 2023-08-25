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

func MakeIntegratedAddress(address string, amount int, comment string, destination int) string {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "MakeIntegratedAddress",
		"params": map[string]interface{}{
			"address": address,
			"payload_rpc": []map[string]interface{}{
				{
					"name":     "C",
					"datatype": "S",
					"value":    comment,
				},
				{
					"name":     "V",
					"datatype": "U",
					"value":    amount,
				},
				{
					"name":     "N",
					"datatype": "U",
					"value":    1,
				}, {
					"name":     "D",
					"datatype": "U",
					"value":    destination,
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return ""
	}

	// Retrieve IP address, wallet port, and node port from environment variables
	ip := os.Getenv("DERO_SERVER_IP")
	derodPort := os.Getenv("DERO_WALLET_PORT")

	// Construct the URL using the retrieved IP address and wallet port
	url := fmt.Sprintf("http://%s:%s/json_rpc", ip, derodPort)

	// Define request for node
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return ""
	}

	// Retrieve authentication credentials from environment variables
	username := os.Getenv("USER")
	password := os.Getenv("PASS")

	// Set basic authentication for the request
	request.SetBasicAuth(username, password)
	request.Header.Set("Content-type", "application/json")
	fmt.Println("\nRequest: ", request)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return ""
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	var mapResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &mapResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return ""
	}

	// Print the entire mapResponse map
	fmt.Println("\nmapResponse:", mapResponse)

	// Print the entire mapResponse map
	var outputMessage string
	for key, value := range mapResponse {
		formattedValue, _ := json.MarshalIndent(value, "", "  ")
		outputMessage += fmt.Sprintf("```%s: %s\n```", key, formattedValue)
	}
	return outputMessage
}
