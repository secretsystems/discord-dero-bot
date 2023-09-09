package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetJSON() (map[string]interface{}, error) {
	url := "https://derostats.io/json"

	// Send a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	// Decode the JSON response into a map
	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
