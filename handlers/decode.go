package handlers

import "github.com/bwmarrin/discordgo"

func handleDecode(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
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
