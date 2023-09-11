package handlers

import "github.com/bwmarrin/discordgo"

func handleRegistration(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
	components := createRegisterModalComponents()
	modal := NewModal(session, interaction, "register_"+interaction.Interaction.Member.User.ID, "Secret Discord Server Registration", components)
	modal.Show()
}

func createRegisterModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "register",
				Label:       "Register DERO Address/Name with the Server",
				Style:       discordgo.TextInputShort,
				Placeholder: "dero1q wallet address or wallet-name",
				Required:    true,
			},
		}},
	}
}
