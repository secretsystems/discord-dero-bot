package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetSC(scid, key string) string {
	// Define JSON struct
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "DERO.GetSC",
		"params": map[string]interface{}{
			"scid":      scid,
			"code":      false,
			"variables": true,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return "Error marshaling JSON"
	}

	// Construct the URL using the retrieved IP address and wallet port
	url := fmt.Sprintf("http://%s:%s/json_rpc", DeroServerIP, DeroServerPort)

	// Define request for node
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "Error creating Request"
	}

	// Set basic authentication for the request
	request.SetBasicAuth(deroUser, deroPass)
	request.Header.Set("Content-Type", "application/json")

	// fmt.Println("\nRequest: ", request)
	client := http.DefaultClient

	// Read response body
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return "Error Senging HTTP Post"
	}
	defer response.Body.Close()

	// Check if the response status code indicates authorization required
	if response.StatusCode == http.StatusUnauthorized {
		log.Printf("Authorization required for the request")
		return "Authorization Required"
	}

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	// Check if the response body is not empty
	if len(responseBody) > 0 {
		var mapResponse map[string]interface{}
		err = json.Unmarshal(responseBody, &mapResponse)
		if err != nil {
			log.Printf("Error decoding response JSON: %v", err)
			return "Error Decoding JSON"
		}

		// Navigate through nested structure
		result, ok := mapResponse["result"].(map[string]interface{})
		if !ok {
			return fmt.Sprintf("Key 'result' not found in the response")
		}

		value, ok := result[key]
		if !ok {
			return fmt.Sprintf("Key '%s' not found in the response", key)
		}

		formattedValue, err := json.MarshalIndent(value, "", " ")
		if err != nil {
			log.Printf("Error formatting JSON value: %v", err)
			return ""
		}

		return string(formattedValue)
	}

	return "Response body is empty"
}
