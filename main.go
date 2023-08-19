package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"fuck_you.com/bot"
)

var (
	BotToken = flag.String("token", os.Getenv("BOT_TOKEN"), "Bot access token")
)

func main() {
	bot, err := bot.NewBot(*BotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	defer bot.Close()

	// Open the bot session
	if err := bot.Open(); err != nil {
		log.Fatalf("Error opening bot session: %v", err)
	}

	// Wait for an interrupt signal to close the bot
	fmt.Println("Bot is running. Press Ctrl+C to stop.")
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	<-channel
}
