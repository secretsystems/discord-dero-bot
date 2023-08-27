package handlers

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Import the godotenv package
)

var (
	chatGptApi string
)

func init() {
	// chatGPT
	chatGptApi = os.Getenv("OPEN_AI_TOKEN")
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
