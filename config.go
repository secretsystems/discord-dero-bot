// config.go
package main

import (
	"log"

	"github.com/joho/godotenv"
)

var (

	// chatGPT
	ChatGptApi string
)

func loadConfig() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Adding setting up configs")
	// Read environment variables

	// Ensure the directory exists

}
