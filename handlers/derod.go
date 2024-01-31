package handlers

import (
	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func HandleGetInfoDerod(session *discordgo.Session, message *discordgo.MessageCreate) {
	endpoint := "127.0.0.1:10102"
	outputMessage := dero.GetInfoDerod(endpoint)
	// Send the entire response to session
	session.ChannelMessageSend(message.ChannelID, "Node Info:\n```\n"+outputMessage+"```")
}
