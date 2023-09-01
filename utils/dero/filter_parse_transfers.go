package dero

import (
	"fmt"
	"log"
	"time"
)

// Helper function to check if the entry has the desired value in "D"
func hasDesiredValue(entry TransferEntry, desiredValue int) bool {
	for _, arg := range entry.PayloadRPC {
		if arg.Name == "V" && arg.Datatype == "U" && arg.Value == desiredValue {
			return true
		}
	}
	return false
}

// FilterAndPrintTransactions is a function that filters parsed transactions with PayloadRPC and prints them to the screen.
func FilterAndPrintTransactions(userID string, pongAmount int) (bool, error) {
	// Fetch transfers and parse the JSON response
	log.Println("Fetching transfers response...")
	responseBody, err := FetchDeroTransfers()
	if err != nil {
		fmt.Printf("Error fetching transfers: %v\n", err)
		return false, err
	}

	// Parse the JSON response
	log.Println("Parsing transfers response...")
	entries, err := ParseTransfersResponse(responseBody)
	if err != nil {
		fmt.Printf("Error parsing transfers response: %v\n", err)
		return false, err
	}

	// Get the current month
	currentMonth := time.Now().Month()

	// Print filtered transactions with PayloadRPC and additional filters
	fmt.Printf("Transactions with PayloadRPC value: %s\n", userID)
	for _, entry := range entries {
		for _, arg := range entry.PayloadRPC {
			if arg.Name == "C" && arg.Datatype == "S" && arg.Value == userID &&
				// hasDesiredValue(entry, pongAmount) &&
				entry.Time.Month() == currentMonth {
				// This transaction meets the criteria, print it to the screen
				fmt.Printf("TxID: %s\n", entry.TxID)
				fmt.Printf("Amount: %d\n", entry.Amount)
				fmt.Printf("Time: %s\n", entry.Time)
				fmt.Println("PayloadRPC:")
				for _, arg := range entry.PayloadRPC {
					fmt.Printf("Name: %s, Datatype: %s, Value: %v\n", arg.Name, arg.Datatype, arg.Value)
				}
				fmt.Println("------")
				return true, nil
			}
		}
	}

	return false, nil
}
