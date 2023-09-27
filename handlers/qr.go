package handlers

import (
	"bytes"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/skip2/go-qrcode"
)

func handleQRModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Command invoked in DM")
		RespondWithMessage(session, interaction, "This command cannot be used in DMs.")
		return
	}
	components := createQRModalComponents()
	modal := NewModal(session, interaction, "qr_"+interaction.Interaction.Member.User.ID, "Create a qr code", components)
	modal.Show()
}

func createQRModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "qr",
				Label:       "Content of the QR Code",
				Style:       discordgo.TextInputParagraph,
				Placeholder: "have fun :)",
				Required:    true,
				MaxLength:   1337,
				MinLength:   1,
			},
		}},
	}
}

func handleQRInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Interaction received in DM")
		RespondWithMessage(session, interaction, "This interaction cannot be processed in DMs.")
		return
	}
	data := interaction.ModalSubmitData()

	// Helper function to get a TextInput value by index
	getTextInputValue := func(index int) string {
		return data.Components[index].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	text := getTextInputValue(0)

	// Generate QR code
	qrCode, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		log.Println("Error generating QR code:", err)
		return
	}

	qrCodeImage, err := qrCode.PNG(256) // Generate QR code image as PNG
	if err != nil {
		log.Printf("Error getting qr code: %v", err)
		return
	}

	// Create a Discord message with the QR code image as an attachment
	msg := &discordgo.MessageSend{
		Content: "Here's your QR code:",
		Files: []*discordgo.File{
			{
				Name:   "qrcode.png",
				Reader: bytes.NewReader(qrCodeImage),
			},
		},
	}

	// Respond with an ephemeral message containing the QR code image
	respondWithImageComponents(session, interaction, "Your Input: "+text, nil, msg)
}

func respondWithImageComponents(session *discordgo.Session, interaction *discordgo.InteractionCreate, message string, components []discordgo.MessageComponent, msg *discordgo.MessageSend) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
			Embeds:     msg.Embeds,
			Files:      msg.Files,
		},
	})
	if err != nil {
		log.Println("Error responding with message and components:", err)
	}
}
