package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func HandleUnregister(session *discordgo.Session, message *discordgo.MessageCreate) {

	content := message.Content
	if content == "!unregister" {
		loadUserMappings()
		// Extract the user ID

		log.Printf("User ID: %s", message.Author.ID) // Added "User ID:"

		userMappingsMutex.Lock()
		defer userMappingsMutex.Unlock()

		// Check if the user is registered
		if _, exists := userMappings[message.Author.ID]; exists {
			delete(userMappings, userID)
			saveUserMappings()

			// Remove the registered role and add the unregistered role
			err := session.GuildMemberRoleRemove(message.GuildID, message.Author.ID, desiredRole)
			if err != nil {
				log.Println("Error removing role from member:", err)
			}

			_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> has been successfully unregistered.", message.Author.ID)) // Added "successfully"
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		} else {
			_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> was not registered.", message.Author.ID))
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}
	}
}
