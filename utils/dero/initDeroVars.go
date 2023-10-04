package dero

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Import the godotenv package
)

var (
	// dero
	DeroServerIP     string
	DeroServerPort   string
	deroWalletPort   string
	deroUser         string
	deroPass         string
	pongAmount       = "1337331"
	pongDir          string
	pongDB           string
	iAddressTextFile string
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//dero
	DeroServerIP = os.Getenv("DERO_SERVER_IP")
	deroWalletPort = os.Getenv("DERO_WALLET_PORT")
	DeroServerPort = os.Getenv("DERO_NODE_PORT")
	deroUser = os.Getenv("USER")
	deroPass = os.Getenv("PASS")
	homeDir, _ := os.UserHomeDir()
	pongDir = homeDir + "/dero-utils"
	pongDB = pongDir + "/" + pongAmount + ".sales.db"
	iAddressTextFile = pongDir + "/" + pongAmount + ".iaddress.txt"

	if _, err := os.Stat(pongDir); os.IsNotExist(err) {
		err := os.Mkdir(pongDir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}
