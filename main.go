package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	loadConfig()

	bot, err := initDiscordBot(BotToken)
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}
	defer bot.Close()

	handleDEROFunctionality()
	initChatGPT()

	fmt.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
}
