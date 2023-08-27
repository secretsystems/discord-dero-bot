package main

import (
	"discord-dero-bot/bot"      // Update with the correct import path for your bot package
	"discord-dero-bot/handlers" // Update with the correct import path for your handlers package
	"discord-dero-bot/utils/dero"
	"log"
	"os"
	"os/signal"
)

func main() {
	loadConfig()

	// Initialize the bot
	bot, err := bot.NewBot(BotToken) // Replace with the actual initialization function
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}

	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening Discord bot connection: %v", err)
	}
	defer bot.Close()

	// Get the Discord session from the bot instance
	discord := bot.GetDiscordSession()

	// Register interaction handlers
	handlers.AddHandlers(discord, AppID, GuildID)
	handlers.AddModals(discord, AppID, GuildID, ResultsChannel)
	handlers.RegisterSlashCommands(discord, AppID, GuildID)

	dero.HandleDEROFunctionality()

	log.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
	handlers.Cleanup(discord, AppID, GuildID)
	log.Println("Bot is cleaning up.")

}
