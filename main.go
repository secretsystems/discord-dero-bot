package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"fuck_you.com/bot" // Import your bot package

	"github.com/joho/godotenv" // Import the godotenv package
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
)

func init() {
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

	// Ensure the directory exists
	if _, err := os.Stat(PongDir); os.IsNotExist(err) {
		err := os.Mkdir(PongDir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}

func printInit() {
	fmt.Println("Initialized variables:")
	fmt.Printf("BotToken: %s\n", BotToken)
	fmt.Printf("GuildID: %s\n", GuildID)
	fmt.Printf("AppID: %s\n", AppID)
	fmt.Printf("DeroServerIP: %s\n", DeroServerIP)
	fmt.Printf("DeroWalletPort: %s\n", DeroWalletPort)
	fmt.Printf("DeroNodePort: %s\n", DeroNodePort)
	fmt.Printf("DeroUser: %s\n", DeroUser)
	fmt.Printf("DeroPass: %s\n", DeroPass)
	fmt.Printf("PongDir: %s\n", PongDir)
	fmt.Printf("PongDB: %s\n", PongDB)
	fmt.Printf("IAddressTextFile: %s\n", IAddressTextFile)
}

func main() {
	printInit()

	// Create a new bot instance using the provided token
	bot, err := bot.NewBot(BotToken)
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

	// Wait for an interrupt signal to close the program
	<-channel
}
