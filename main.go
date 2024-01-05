package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/secretsystems/discord-dero-bot/bot"      // Update with the correct import path for your bot package
	"github.com/secretsystems/discord-dero-bot/handlers" // Update with the correct import path for your handlers package
	"github.com/secretsystems/discord-dero-bot/utils/dero"
	"github.com/secretsystems/discord-dero-bot/utils/monero"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	// discord
	botToken string
	guildID  string
	appID    string
)

func init() {

	// 	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	//discord
	botToken = os.Getenv("BOT_TOKEN")
	guildID = os.Getenv("GUILD_ID")
	appID = os.Getenv("APP_ID")

	log.Printf("Initializing DERO\n")
	// Call FetchAndParseTransfers function from the utils package
	transferEntries, err := dero.FetchAndParseTransfers()

	if err != nil {
		log.Printf("Error fetching and parsing transfers: %v", err)
	} else {
		// Process the fetched and parsed transfer entries
		log.Printf("Fetched and parsed %d transfer entries.\n", len(transferEntries))
	}
	dero.GetDeroWalletBalance()
	monero.GetWalletBalance()
}

func main() {
	// Initialize the bot

	bot, err := bot.NewBot(botToken) // Replace with the actual initialization function
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening Discord bot connection: %v", err)
	}
	defer bot.Close()

	// Get the Discord session from the bot instance
	session := bot.GetDiscordSession()

	// Register interaction handlers

	handlers.AddHandlers(session, appID)
	handlers.AddModals(session, appID)
	registrationBucket := handlers.NewTokenBucket(1, 1, time.Second*4) // Allow 1 request every 5 seconds
	handlers.RegisterSlashCommands(session, appID, guildID, registrationBucket)

	bot.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
	cleanupBucket := handlers.NewTokenBucket(1, 1, time.Second*4) // Allow 1 request every 5 seconds
	handlers.Cleanup(session, appID, guildID, cleanupBucket)
	log.Println("Bot is cleaning up.")

}
