package dero

import "fmt"

func FetchAndParseTransfers() ([]TransferEntry, error) {
	// Call the FetchDeroTransfers function to obtain the JSON response
	responseBody, err := FetchDeroTransfers()
	if err != nil {
		fmt.Printf("Error fetching transfers: %v\n", err)
		return nil, err
	}

	// Parse the JSON response and extract the "height" values
	fmt.Println("Parsing transfers response...")
	entries, err := ParseTransfersResponse(responseBody)
	if err != nil {
		fmt.Printf("Error parsing transfers response: %v\n", err)
		return nil, err
	}

	return entries, nil
}
