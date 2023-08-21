// utils/parse_transfers.go

package utils

import (
	"encoding/json"
	"fmt"
)

// Define a struct to match the JSON entry structure
type TransferEntry struct {
	Height uint64 `json:"height"`
	// Add other fields here if needed
}

func ParseTransfersResponse(responseBody []byte) ([]TransferEntry, error) {
	var responseData map[string]interface{}

	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	entriesRaw := responseData["result"].(map[string]interface{})["entries"].([]interface{})

	var entries []TransferEntry

	for _, entryRaw := range entriesRaw {
		entryData, err := json.Marshal(entryRaw)
		if err != nil {
			return nil, fmt.Errorf("error marshaling entry: %v", err)
		}

		var entry TransferEntry
		if err := json.Unmarshal(entryData, &entry); err != nil {
			return nil, fmt.Errorf("error unmarshaling entry: %v", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
