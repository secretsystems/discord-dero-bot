package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleUnregister(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	if content == "!unregister" {
		// Extract the user ID
		userID := message.Author.ID

		userMappingsMutex.Lock()
		defer userMappingsMutex.Unlock()

		// Check if the user is registered
		if _, exists := userMappings[userID]; exists {
			delete(userMappings, userID)
			saveUserMappings()

			_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> has been unregistered.", userID))
			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
			}
		} else {
			_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> was not registered.", userID))
			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
			}
		}
	}
}
