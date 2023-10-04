package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SplitIntegratedAddress(userInput string) string {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "SplitIntegratedAddress",
		"params": map[string]string{
			"integrated_address": userInput,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return "Error marshaling JSON"
	}

	// Construct the URL using the retrieved IP address and wallet port
	url := fmt.Sprintf("http://%s:%s/json_rpc", DeroServerIP, deroWalletPort)

	// Define request for node
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "Error creating request"
	}

	// Set basic authentication for the request
	request.SetBasicAuth(deroUser, deroPass)
	request.Header.Set("Content-type", "application/json")
	// fmt.Println("\nRequest: ", request)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return "Error sending HTTP Post"
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	// log.Printf("Response Body: %v", string(responseBody))

	var mapResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &mapResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return "Error decoding response JSON"
	}

	// Print the entire mapResponse map
	// fmt.Println("\nmapResponse:", mapResponse)

	// Print the entire mapResponse map
	var outputMessage string
	for key, value := range mapResponse {
		formattedValue, _ := json.MarshalIndent(value, "", "  ")
		outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
	}
	return outputMessage
}
