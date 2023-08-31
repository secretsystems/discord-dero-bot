package main

import (
	"discord-dero-bot/bot"      // Update with the correct import path for your bot package
	"discord-dero-bot/handlers" // Update with the correct import path for your handlers package
	"discord-dero-bot/utils/dero"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	// discord
	botToken       = flag.String("token", "", "Bot access token")
	guildID        = flag.String("guild", "", "Test guild ID")
	appID          = flag.String("app", "", "Application ID")
	resultsChannel = flag.String("results", "", "Channel where send survey results to")
)

func init() {
	flag.Parse()
}

func main() {
	// Initialize the bot

	bot, err := bot.NewBot(*botToken) // Replace with the actual initialization function
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}

	log.Printf("Initializing DERO\n")
	// Call FetchAndParseTransfers function from the utils package
	transferEntries, err := dero.FetchAndParseTransfers()

	if err != nil {
		log.Printf("Error fetching and parsing transfers: %v", err)
	} else {
		// Process the fetched and parsed transfer entries
		log.Printf("Fetched and parsed %d transfer entries.\n", len(transferEntries))
	}

	bot.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening Discord bot connection: %v", err)
	}
	defer bot.Close()

	// Get the Discord session from the bot instance
	session := bot.GetDiscordSession()

	// Register interaction handlers

	handlers.AddHandlers(session, *appID, *guildID)
	handlers.AddModals(session, *appID, *guildID, *resultsChannel)
	handlers.RegisterSlashCommands(session, *appID, *guildID)

	log.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
	handlers.Cleanup(session, *appID, *guildID)
	log.Println("Bot is cleaning up.")

}
