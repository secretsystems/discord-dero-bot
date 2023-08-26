package handlers

import (
	"discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func HandleGetInfoDerod(discord *discordgo.Session, message *discordgo.MessageCreate) {

	outputMessage := dero.GetInfo()
	// Send the entire response to Discord
	discord.ChannelMessageSend(message.ChannelID, "Node Info:\n```\n"+outputMessage+"```")
}
