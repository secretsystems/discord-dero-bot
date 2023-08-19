package bot

import (
	"strings"

	"fuck_you.com/handlers"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	discord *discordgo.Session
}

func NewBot(token string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		discord: discord,
	}

	discord.AddHandler(bot.newMessage)
	return bot, nil
}

func (bot *Bot) Open() error {
	return bot.discord.Open()
}

func (bot *Bot) Close() {
	bot.discord.Close()
}

func (bot *Bot) newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "!compliment"):
		bot.discord.ChannelMessageSend(message.ChannelID, "You are a wonderful human!")
	case strings.Contains(message.Content, "!insult"):
		bot.discord.ChannelMessageSend(message.ChannelID, "We don't say that stuff around here!")
	case strings.Contains(message.Content, "!decode"):
		handlers.HandleIntegratedAddress(bot.discord, message)
	}
}
