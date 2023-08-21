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

type WalletInfo struct {
	Address      string
	IsRegistered bool
}

func WalletNameToAddress(walletName string) string {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "DERO.NameToAddress",
		"params": map[string]string{
			"name": walletName,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return ""
	}

	ip := os.Getenv("DERO_SERVER_IP")
	derodPort := os.Getenv("DERO_NODE_PORT")
	url := fmt.Sprintf("http://%s:%s/json_rpc", ip, derodPort)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return ""
	}

	username := os.Getenv("USER")
	password := os.Getenv("PASS")

	request.SetBasicAuth(username, password)
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return ""
	}

	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	if len(responseBody) > 0 {
		var mapResponse map[string]interface{}
		err = json.Unmarshal(responseBody, &mapResponse)
		if err != nil {
			log.Printf("Error decoding response JSON: %v", err)
			return ""
		}

		result, resultExists := mapResponse["result"].(map[string]interface{})
		if resultExists {
			addr, addrExists := result["address"].(string)
			if addrExists {
				return addr
			}
		}
	}

	return ""
}