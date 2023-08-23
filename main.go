package main

import (
	"log"
	"os"
	"os/signal"

	"fuck_you.com/bot"      // Update with the correct import path for your bot package
	"fuck_you.com/handlers" // Update with the correct import path for your handlers package
)

func main() {
	loadConfig()

	// Initialize the bot
	botInstance, err := bot.NewBot(BotToken) // Replace with the actual initialization function
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}

	err = botInstance.Open()
	if err != nil {
		log.Fatalf("Error opening Discord bot connection: %v", err)
	}
	defer botInstance.Close()

	// Get the Discord session from the bot instance
	discordSession := botInstance.GetDiscordSession()

	// Register interaction handlers
	handlers.RegisterInteractionHandlersFromHandlers(discordSession, AppID, GuildID)

	handleDEROFunctionality()
	initChatGPT()

	log.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
}
