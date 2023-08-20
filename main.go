package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"fuck_you.com/bot" // Import your bot package
	// Import your utilities package
	"github.com/joho/godotenv" // Import the godotenv package
)

var (
	BotToken string
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a new bot instance using the provided token
	bot, err := bot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Close the bot when the main function exits
	defer bot.Close()

	// Open the bot session
	if err := bot.Open(); err != nil {
		log.Fatalf("Error opening bot session: %v", err)
	}

	// Print a message indicating that the bot is running
	fmt.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Start the transfers fetching routine
	// go runTransfersFetching()

	// Wait for an interrupt signal to close the program
	<-channel
}

// func runTransfersFetching() {
// 	// Fetch Dero transfers in a loop
// 	for {
// 		// Fetch transfers data using the utils package
// 		transfersData, err := utils.FetchDeroTransfers()
// 		if err != nil {
// 			log.Printf("Error fetching Dero transfers: %v", err)
// 		} else {
// 			// Format the JSON data for pretty printing
// 			prettyJSON := utils.FormatJSON(transfersData)
// 			// fmt.Println(prettyJSON)
// 		}

// 		// Sleep for a while before fetching transfers again
// 		time.Sleep(18 * time.Second) // Adjust the sleep duration as needed
// 	}
// }
