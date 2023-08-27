package monero

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// monero
	moneroServerIP   string
	moneroServerPort string
	moneroWalletPort string
	moneroUser       string
	moneroPass       string
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// monero
	moneroServerIP = os.Getenv("MONERO_SERVER_IP")
	moneroServerPort = os.Getenv("MONERO_SERVER_PORT")
	moneroWalletPort = os.Getenv("MONERO_WALLET_PORT")
	moneroUser = os.Getenv("MONERO_WALLET_USER")
	moneroPass = os.Getenv("MONERO_WALLET_PASS")
}
