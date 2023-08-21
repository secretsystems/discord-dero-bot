// utils/dero/fetch_and_parse_transfers.go

package utils

func FetchAndParseTransfers() ([]TransferEntry, error) {
	// Call the FetchDeroTransfers function to obtain the JSON response
	responseBody, err := FetchDeroTransfers()
	if err != nil {
		return nil, err
	}

	// Parse the JSON response and extract the "height" values
	entries, err := ParseTransfersResponse(responseBody)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
