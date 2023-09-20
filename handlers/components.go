package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Modal struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	CustomID    string
	Title       string
	Components  []discordgo.MessageComponent
}

// Commands represents the application commands.
var Commands = []discordgo.ApplicationCommand{
	{
		Name:        "encode",
		Description: "Encode Integrated Address",
	},
	{
		Name:        "trade-dero-xmr",
		Description: "Trade DERO-XMR",
	},
	{
		Name:        "decode",
		Description: "Decode Integrated Address",
	},
	{
		Name:        "giftbox",
		Description: "Get a DERO gift box!",
	},
	{
		Name:        "register",
		Description: "Register your DERO wallet address/name with the server!",
	},
	{
		Name:        "buy-dero-with-crypto",
		Description: "Purchase DERO from the `secret-wallet` with crypto",
	},
	{
		Name:        "buy-dero-with-fiat",
		Description: "Purchase DERO from the `secret-wallet` with FIAT",
	},
}

// DefineHandlers defines the component and command handlers.
func DefineHandlers(session *discordgo.Session, appID string) map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	handlers := make(map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string))

	// Add your command handlers
	handlers["encode"] = handleEncode
	handlers["trade-dero-xmr"] = handleTradeDeroXmr
	handlers["decode"] = handleDecode
	handlers["giftbox"] = handleGiftbox
	handlers["register"] = handleRegistration

	// Add your component handlers
	handlers["fd_yes"] = handleFdYes
	handlers["fd_no"] = handleFdNo

	return handlers
}

// RespondWithMessage sends a response message to the interaction.
func RespondWithMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate, message string) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}

// RespondWithModal sends a modal response to the interaction.
func RespondWithModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, customID, title string, components []discordgo.MessageComponent) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   customID,
			Title:      title,
			Components: components,
		},
	})
	if err != nil {
		log.Printf("Response Error: %v", err)
	}
}

func respondWithMessageAndComponents(session *discordgo.Session, interaction *discordgo.InteractionCreate, message string, components []discordgo.MessageComponent) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	})
	if err != nil {
		panic(err)
	}
}

// NewModal creates a new Modal instance.
func NewModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, customID, title string, components []discordgo.MessageComponent) *Modal {
	return &Modal{
		Session:     session,
		Interaction: interaction,
		CustomID:    customID,
		Title:       title,
		Components:  components,
	}
}
func (m *Modal) Show() {
	err := m.Session.InteractionRespond(m.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   m.CustomID,
			Title:      m.Title,
			Components: m.Components,
		},
	})
	if err != nil {
		panic(err)
	}
}
