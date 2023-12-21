package dero

import (
	"fmt"
	"log"
)

func FetchAndParseTransfers() ([]TransferEntry, error) {
	// Call the FetchDeroTransfers function to obtain the JSON response
	log.Println("Fetching transfers response...")
	responseBody, err := FetchDeroTransfers()
	if err != nil {
		fmt.Printf("Error fetching transfers: %v\n", err)
		return nil, err
	}

	// Parse the JSON response and extract the "height" values
	log.Println("Parsing transfers response...")
	entries, err := ParseTransfersResponse(responseBody)
	if err != nil {
		fmt.Printf("Error parsing transfers response: %v\n", err)
		return nil, err
	}

	// Filter transactions with a PayloadRPC
	filteredEntries := []TransferEntry{}
	for _, entry := range entries {
		if len(entry.PayloadRPC) > 0 {
			// This transaction has a PayloadRPC, include it in the filtered list
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries, nil
}
