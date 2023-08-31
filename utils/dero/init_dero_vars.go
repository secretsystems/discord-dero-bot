package dero

import (
	"log"
	"os"
)

var (
	// dero
	deroServerIP     string
	deroServerPort   string
	deroWalletPort   string
	deroUser         string
	deroPass         string
	pongAmount       = "1337331"
	pongDir          string
	pongDB           string
	iAddressTextFile string
)

func init() {

	//dero
	deroServerIP = "192.168.12.208"
	deroWalletPort = "10103"
	deroServerPort = "10102"
	deroUser = "user"
	deroPass = "pass"
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
