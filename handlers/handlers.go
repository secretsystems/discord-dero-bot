package handlers

import (
	"log"

	"github.com/joho/godotenv" // Import the godotenv package
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
