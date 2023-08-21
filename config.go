// config.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BotToken         string
	GuildID          string
	AppID            string
	DeroServerIP     string
	DeroWalletPort   string
	DeroUser         string
	DeroPass         string
	PongAmount       = "1337331"
	PongDir          string
	PongDB           string
	IAddressTextFile string
	DeroNodePort     string
	ChatGptApi       string
)

func loadConfig() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	BotToken = os.Getenv("BOT_TOKEN")
	GuildID = os.Getenv("GUILD_ID")
	AppID = os.Getenv("APP_ID")
	DeroServerIP = os.Getenv("DERO_SERVER_IP")
	DeroWalletPort = os.Getenv("DERO_WALLET_PORT")
	DeroUser = os.Getenv("USER")
	DeroPass = os.Getenv("PASS")
	homeDir, _ := os.UserHomeDir()
	PongDir = homeDir + "/dero-utils"
	PongDB = PongDir + "/" + PongAmount + ".sales.db"
	IAddressTextFile = PongDir + "/" + PongAmount + ".iaddress.txt"
	DeroNodePort = os.Getenv("DERO_NODE_PORT")
	ChatGptApi = os.Getenv("OPEN_AI_TOKEN")

	// Ensure the directory exists
	if _, err := os.Stat(PongDir); os.IsNotExist(err) {
		err := os.Mkdir(PongDir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}
