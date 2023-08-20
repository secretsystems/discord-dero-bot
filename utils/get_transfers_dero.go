// utils/gettransfers_dero.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchDeroTransfers() ([]byte, error) {
	url := "http://192.168.12.208:10103/json_rpc"
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
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("user", "pass")
	req.Header.Set("Content-type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FormatJSON(jsonData []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, jsonData, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return ""
	}
	return prettyJSON.String()
}
