package bot

import (
	"strings"

	"fuck_you.com/handlers"
	"github.com/bwmarrin/discordgo"
)

var CommandHandlers = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!compliment": handlers.HandleCompliment,
	"!insult":     handlers.HandleInsult,
	"!decode":     handlers.HandleIntegratedAddress,
	"!lookup":     handlers.HandleWalletName,
	"!derod":      handlers.HandleGetInfoDerod,
	"!monerod":    handlers.HandleGetInfoMonerod,
	"!quote":      handlers.HandleQuoteRequest,
	"!markets":    handlers.HandleMarketsRequest,
	"!help":       handlers.HandleHelp,
	// "!membership"  handlers.HandleMembership
}

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

	for command, handler := range CommandHandlers {
		if strings.Contains(message.Content, command) {
			handler(discord, message)
			return
		}
	}
}
