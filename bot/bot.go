package bot

import (
	"log"
	"reflect"
	"strings"

	"discord-dero-bot/handlers"

	"github.com/bwmarrin/discordgo"
)

var PingHandlers = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"!compliment": handlers.HandleMessage,
	"!insult":     handlers.HandleMessage,
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
	"!shop":       handlers.HandleShop,
	// "!membership"  handlers.HandleMembership
}

type Bot struct {
	DiscordSession *discordgo.Session // Exported field for accessing the session session
}

func NewBot(BotToken string) (*Bot, error) {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		DiscordSession: session, // Store the session in the Bot instance
	}

	bot.AddHandler(bot.NewMessage)
	bot.AddHandler(bot.OnReady) // Add this line to register the onReady handler
	bot.AddHandler(bot.OnGeneric)
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

func (bot *Bot) AddHandler(handler interface{}) func() {
	return bot.DiscordSession.AddHandler(handler)
}

func (bot *Bot) OnGeneric(s *discordgo.Session, event interface{}) {
	t := reflect.TypeOf(event)

	log.Printf("GENERIC EVENT: %v\n", t)
	switch cast := event.(type) {
	case *discordgo.TypingStart:
		log.Printf("TYPING EVENT: User %v is typing on channel %v\n", cast.UserID, cast.ChannelID)
	case *discordgo.MessageCreate:
		log.Printf("MESSAGE EVENT: %v | %v", cast.Author, cast.Content)
	}
}

func (bot *Bot) OnReady(session *discordgo.Session, ready *discordgo.Ready) {

}

func (bot *Bot) NewMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == bot.DiscordSession.State.User.ID {
		return
	}
	for command, handler := range PingHandlers {
		if strings.HasPrefix(message.Content, command) {
			handler(bot.DiscordSession, message)
			return
		}
	}
}
