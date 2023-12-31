package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
)

// Handle !grok command
func HandleGrok(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, dero.GetGrok())
}
