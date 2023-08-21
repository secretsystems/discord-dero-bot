package handlers

import (
	"strings"

	"fuck_you.com/utils/dero" // Import the dero package from your project
	"github.com/bwmarrin/discordgo"
)

func HandleTip(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if strings.HasPrefix(content, "!tip ") {
		address := strings.TrimPrefix(content, "!tip ")
		dero.MakeTransfer(address) // Call MakeTransfer() with the extracted address
		discord.ChannelMessageSend(message.ChannelID, "Tip sent!")
	}
}
