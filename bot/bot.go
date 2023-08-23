package bot

import (
	"strings"

	"fuck_you.com/handlers"
	"github.com/bwmarrin/discordgo"
)

var CommandHandlers = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!compliment": handlers.HandleMessage,
	"!insult":     handlers.HandleMessage,
	"!decode":     handlers.HandleIntegratedAddress,
	"!lookup":     handlers.HandleWalletName,
	"!derod":      handlers.HandleGetInfoDerod,
	"!monerod":    handlers.HandleGetInfoMonerod,
	"!quote":      handlers.HandleQuoteRequest,
	"!markets":    handlers.HandleMarketsRequest,
	"!help":       handlers.HandleHelp,
	"!bot":        handlers.HandleChat,
	"!tip":        handlers.HandleTip,
	"!register":   handlers.HandleRegister,
	"!unregister": handlers.HandleUnregister,
	// "!membership"  handlers.HandleMembership
}

type Bot struct {
	DiscordSession *discordgo.Session // Exported field for accessing the Discord session
}

func NewBot(BotToken string) (*Bot, error) {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		DiscordSession: discord, // Store the session in the Bot instance
	}

	discord.AddHandler(bot.NewMessage)
	discord.AddHandler(bot.OnReady) // Add this line to register the onReady handler
	return bot, nil
}

func (bot *Bot) Open() error {
	return bot.DiscordSession.Open()
}

func (bot *Bot) Close() {
	bot.DiscordSession.Close()
}

func (bot *Bot) GetDiscordSession() *discordgo.Session {
	return bot.DiscordSession
}

func (bot *Bot) OnReady(discord *discordgo.Session, ready *discordgo.Ready) {

}

func (bot *Bot) NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == bot.DiscordSession.State.User.ID {
		return
	}
	for command, handler := range CommandHandlers {
		if strings.HasPrefix(message.Content, command) {
			handler(bot.DiscordSession, message)
			return
		}
	}
}
