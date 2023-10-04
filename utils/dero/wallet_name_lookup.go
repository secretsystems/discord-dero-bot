package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WalletInfo struct {
	Address      string
	IsRegistered bool
}

func WalletNameToAddress(input string) (string, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "DERO.NameToAddress",
		"params": map[string]string{
			"name": input,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return "", err
	}

	url := fmt.Sprintf("http://%s:%s/json_rpc", DeroServerIP, DeroServerPort)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	request.SetBasicAuth(deroUser, deroPass)
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return "", err
	}

	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	if len(responseBody) > 0 {
		var mapResponse map[string]interface{}
		err = json.Unmarshal(responseBody, &mapResponse)
		if err != nil {
			log.Printf("Error decoding response JSON: %v", err)
			return "", err
		}

		errorObj, errorExists := mapResponse["error"].(map[string]interface{})
		if errorExists {
			errorMessage, messageExists := errorObj["message"].(string)
			if messageExists {
				return "", fmt.Errorf("DERO API Error: %s", errorMessage)
			}
		}

		result, resultExists := mapResponse["result"].(map[string]interface{})
		if resultExists {
			addr, addrExists := result["address"].(string)
			if addrExists {
				return addr, nil
			}
		}
	}

	return "", fmt.Errorf("Invalid wallet name: %s", input)
}
