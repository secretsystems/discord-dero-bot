package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
)

func HandleSCIDLookup(session *discordgo.Session, message *discordgo.MessageCreate) {
	loadUserMappings()
	content := message.Content
	// fmt.Println("CONTENT: %s", content)

	if !strings.HasPrefix(content, "!scid ") {
		userMessage := "To lookup a DERO scid, use the format: `!scid <scid> <key>`"
		session.ChannelMessageSend(message.ChannelID, userMessage)
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!scid ")
	inputParts := strings.Fields(userInput)

	if len(inputParts) < 2 {
		session.ChannelMessageSend(message.ChannelID, "Invalid input format. Please provide both SCID and key.")
		return
	}

	deroScid := inputParts[0]
	key := inputParts[1]

	result := dero.GetStringKey(deroScid, key)

	switch v := result.(type) {
	case string:
		session.ChannelMessageSend(message.ChannelID, "DERO scid: ```"+v+"```")
	case float64:
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("DERO scid: %g", v))
	default:
		session.ChannelMessageSend(message.ChannelID, "Unexpected result type. Unable to process.")
	}
}
