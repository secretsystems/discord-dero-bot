package handlers

import (
	"discord-dero-bot/utils/monero"

	"github.com/bwmarrin/discordgo"
)

func HandleGetInfoMonerod(session *discordgo.Session, message *discordgo.MessageCreate) {

	outputMessage := monero.GetInfoMonerod()
	// Send the entire response to Discord
	session.ChannelMessageSend(message.ChannelID, "Node Info:\n```\n"+outputMessage+"```")

}
