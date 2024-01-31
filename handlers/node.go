package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
)

func handleNodeModal(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	appID string,
) {
	components := createNodeModalComponents()
	modal := NewModal(
		session,
		interaction,
		"node_"+interaction.Interaction.Member.User.ID,
		"Look up node info",
		components,
	)
	modal.Show()
}

func createNodeModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "node",
					Label:       "Look up info on node",
					Style:       discordgo.TextInputShort,
					Placeholder: "89.38.99.117:10102",
					Required:    true,
				},
			},
		},
	}
}

func handleNodeInteraction(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) {
	data := interaction.ModalSubmitData()
	// Helper function to get a TextInput value by index
	getTextInputValue := func(index int) string {
		return data.Components[index].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	text := getTextInputValue(0)
	message := dero.GetInfoDerod(text)
	message = "Node Info:\n```\n" + message + "```"
	RespondWithMessage(
		session,
		interaction,
		message,
	)
}
