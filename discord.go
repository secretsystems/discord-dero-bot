// discord.go
package main

import (
	"fmt"
	"log"

	"fuck_you.com/bot"
)

func initDiscordBot(BotToken string) (*bot.Bot, error) {
	// Create a new bot instance using the provided token
	bot, err := bot.NewBot(BotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Open the bot session
	if err := bot.Open(); err != nil {
		log.Fatalf("Error opening bot session: %v", err)
	}

	fmt.Println("Discord bot initialized.")
	return bot, nil
}
