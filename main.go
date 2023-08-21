package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"fuck_you.com/bot" // Import your bot package
	"fuck_you.com/utils"

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
	go FetchAndPrintTransfers()

	// Wait for an interrupt signal to close the program
	<-channel
}

func FetchAndPrintTransfers() {
	// Call the FetchDeroTransfers function to obtain the JSON response
	responseBody, err := utils.FetchDeroTransfers()
	if err != nil {
		log.Fatalf("Error fetching transfers: %v", err)
	}

	log.Printf("Raw response: %s", responseBody) // Print raw response for debugging

	// Parse the JSON response and extract the "height" values
	entries, err := utils.ParseTransfersResponse(responseBody)
	if err != nil {
		log.Fatalf("Error parsing transfers response: %v", err)
	}

	// Print the "height" values
	for _, entry := range entries {
		fmt.Printf("Height: %d\n", entry.Height)
	}
}
