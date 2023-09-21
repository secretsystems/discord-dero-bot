package bot

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"discord-dero-bot/handlers"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
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
	"!derostats":  handlers.HandleDerostats,
	"!unregister": handlers.HandleUnregister,
	"!shop":       handlers.HandleShop,
	"!music":      handlers.HandleMusic,
	// "!membership"  handlers.HandleMembership
}

type Bot struct {
	DiscordSession  *discordgo.Session // Exported field for accessing the session session
	VoiceConnection *discordgo.VoiceConnection
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
func (bot *Bot) PlayAudio(guildID, voiceChannelID string) error {

	// Join the voice channel
	voiceConnection, err := bot.DiscordSession.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
	if err != nil {
		return err
	}

	bot.VoiceConnection = voiceConnection

	// Stream the audio
	err = bot.StreamAudio("https://somafm.com/thetrip.pls")
	if err != nil {
		return err
	}

	// Leave the voice channel once audio is finished
	bot.VoiceConnection.Disconnect()
	return nil
}

func (bot *Bot) StreamAudio(audioURL string) error {
	// Fetch the audio stream
	resp, err := http.Get(audioURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Encode the audio and play it
	opts := dca.StdEncodeOptions
	encodingSession, err := dca.EncodeMem(resp.Body, opts)
	if err != nil {
		return err
	}

	// Send the audio to Discord
	for {
		select {
		case <-time.After(time.Second): // Adjust the timeout as needed
			return nil
		default:
			// Read audio data and send it to Discord
			buf := make([]byte, 96)
			n, err := encodingSession.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				return nil
			}
			bot.VoiceConnection.OpusSend <- buf[:n]
		}
	}
}
