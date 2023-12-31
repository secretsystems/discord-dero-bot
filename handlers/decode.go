package handlers

import (
	"log"

	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func handleDecodeModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Command invoked in DM")
		RespondWithMessage(session, interaction, "This command cannot be used in DMs.")
		return
	}
	components := createDecodeModalComponents()

	modal := NewModal(session, interaction, "decode_"+interaction.Interaction.Member.User.ID, "Decode DERO Integrated Address", components)
	modal.Show()
}
func createDecodeModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "integrated_address",
				Label:       "Integrated Address",
				Style:       discordgo.TextInputShort,
				Placeholder: "integrated wallet address",
				Required:    true,
			},
		}},
	}
}

func handleDecodeInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Interaction received in DM")
		RespondWithMessage(session, interaction, "This interaction cannot be processed in DMs.")
		return
	}
	data := interaction.ModalSubmitData()
	address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	splitIntegratedAddress := dero.SplitIntegratedAddress(address)
	RespondWithMessage(session, interaction, "```"+splitIntegratedAddress+"```")
}
