package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", os.Getenv("GUILD_ID"), "Test guild ID")
	BotToken = flag.String("token", os.Getenv("BOT_TOKEN"), "Bot access token")
	AppID    = flag.String("app", os.Getenv("APP_ID"), "Application ID")
)

var discord *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	discord, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore bot messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Respond to messages
	switch {
	case strings.Contains(message.Content, "!compliment"):
		discord.ChannelMessageSend(message.ChannelID, "You are a wonderful human!")
	case strings.Contains(message.Content, "!insult"):
		discord.ChannelMessageSend(message.ChannelID, "We don't say that stuff around here!")
	}
}

func run() {
	// Add event handler
	discord.AddHandler(newMessage)

	// Open session
	discord.Open()
	defer discord.Close()

	//Run until code is terminated
	fmt.Println("Bot is kicking ass and takin' names")
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	<-channel
}
func main() {
	run()
}
