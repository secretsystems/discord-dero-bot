package handlers

import (
	"discord-dero-bot/utils/dero"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleIntegratedAddress(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if content == "!decode" {
		discord.ChannelMessageSend(message.ChannelID, "To decode an integrated address: `!decode <integrated_address>`")
		return
	} else if strings.HasPrefix(content, "!decode ") {
		userInput := strings.TrimPrefix(message.Content, "!decode ")
		outputMessage := dero.SplitIntegratedAddress(userInput)

		// Get or create the DM channel for the user
		dmChannel, err := discord.UserChannelCreate(message.Author.ID)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error creating DM channel.")
			return
		}

		// Send the response as a private message
		_, err = discord.ChannelMessageSend(dmChannel.ID, "Integrated Address Response:\n```\n"+outputMessage+"```")
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Error sending DM.")
		}
	}
}
