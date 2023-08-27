// config.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// discord
	BotToken       string
	GuildID        string
	AppID          string
	ResultsChannel string

	// dero
	DeroServerIP     string
	DeroNodePort     string
	DeroWalletPort   string
	DeroUser         string
	DeroPass         string
	PongAmount       = "1337331"
	PongDir          string
	PongDB           string
	IAddressTextFile string

	// monero
	MoneroServerIP   string
	MoneroServerPort string
	MoneroWalletPort string
	MoneroUser       string
	MoneroPass       string

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

	// discord
	BotToken = os.Getenv("BOT_TOKEN")
	GuildID = os.Getenv("GUILD_ID")
	AppID = os.Getenv("APP_ID")
	ResultsChannel = os.Getenv("RESULTS_CHANNEL")

	//dero
	DeroServerIP = os.Getenv("DERO_SERVER_IP")
	DeroWalletPort = os.Getenv("DERO_WALLET_PORT")
	DeroUser = os.Getenv("USER")
	DeroPass = os.Getenv("PASS")
	homeDir, _ := os.UserHomeDir()
	PongDir = homeDir + "/dero-utils"
	PongDB = PongDir + "/" + PongAmount + ".sales.db"
	IAddressTextFile = PongDir + "/" + PongAmount + ".iaddress.txt"
	DeroNodePort = os.Getenv("DERO_NODE_PORT")

	// monero
	MoneroServerIP = os.Getenv("MONERO_SERVER_IP")
	MoneroServerPort = os.Getenv("MONERO_SERVER_PORT")
	MoneroWalletPort = os.Getenv("MONERO_WALLET_PORT")
	MoneroUser = os.Getenv("MONERO_WALLET_USER")
	MoneroPass = os.Getenv("MONERO_WALLET_PASS")

	// chatGPT
	ChatGptApi = os.Getenv("OPEN_AI_TOKEN")

	// Ensure the directory exists
	if _, err := os.Stat(PongDir); os.IsNotExist(err) {
		err := os.Mkdir(PongDir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}
