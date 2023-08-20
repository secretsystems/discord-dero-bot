package utils

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// // DeroServer represents the Dero server details
// type DeroServer struct {
// 	URL      string
// 	Username string
// 	Password string
// }

// // NewDeroServer creates a new DeroServer instance
// func NewDeroServer(url, username, password string) *DeroServer {
// 	return &DeroServer{
// 		URL:      url,
// 		Username: username,
// 		Password: password,
// 	}
// }

// // FetchTransfers fetches Dero transfers from the server
// func (ds *DeroServer) FetchTransfers() ([]byte, error) {
// 	data := map[string]interface{}{
// 		"jsonrpc": "2.0",
// 		"id":      "1",
// 		"method":  "GetTransfer",
// 		"params": map[string]interface{}{
// 			"coinbase": true,
// 			"out":      true,
// 			"in":       true,
// 		},
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", ds.URL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.SetBasicAuth(ds.Username, ds.Password)
// 	req.Header.Set("Content-type", "application/json")

// 	// Print the request going to the URL
// 	fmt.Println("Sending request to:", ds.URL)
// 	fmt.Println("Request body:", string(jsonData))

// 	client := http.DefaultClient
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Print the response received from the server
// 	fmt.Println("Response status:", resp.Status)
// 	responseBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Response body:", string(responseBody))

// 	return responseBody, nil
// }

// // FormatJSON formats JSON data
// func FormatJSON(jsonData []byte) string {
// 	var prettyJSON bytes.Buffer
// 	err := json.Indent(&prettyJSON, jsonData, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error formatting JSON:", err)
// 		return ""
// 	}
// 	return prettyJSON.String()
// }
