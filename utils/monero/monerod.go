package monero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetInfoMonerod() string { // Define JSON struct
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_info",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return "Error marshaling JSON"
	}

	url := fmt.Sprintf("http://%s:%s/json_rpc", moneroServerIP, moneroServerPort)

	// Define request for node
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error marchaling JSNON: %v", err)
		return "Error marshaling JSON"
	}

	// Define request authentication for node
	request.SetBasicAuth(moneroUser, moneroPass)
	request.Header.Set("content-type", "application/json")
	// fmt.Println("\nRequest: ", request)
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return "Error sending HTTP"
	}

	// close out authenticated response
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	// log.Printf("Response Body: %v", string(responseBody))

	var mapResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &mapResponse)
	if err != nil {
		log.Printf("Error decoding resopnse JSON: %v", err)
		return "Error decoding resopnse JSON"
	}

	// Print the entire httpResponse map
	// log.Printf("\nResponse Body: %v", string(responseBody))

	var outputMessage string
	for key, value := range mapResponse {
		formattedValue, _ := json.MarshalIndent(value, "", " ")
		outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
	}
	return outputMessage
}
