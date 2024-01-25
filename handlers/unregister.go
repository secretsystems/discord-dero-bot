package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/exports"
)

func HandleUnregister(session *discordgo.Session, message *discordgo.MessageCreate) {

	content := message.Content
	if content == "!unregister" {
		loadUserMappings()
		// Extract the user ID
		userID := message.Author.ID
		log.Printf("User ID: %s", userID)

		userMappingsMutex.Lock()
		defer userMappingsMutex.Unlock()

		// Check if the user is registered
		if _, exists := userMappings[userID]; exists {
			delete(userMappings, userID)
			saveUserMappings()

			if IsMemberInGuild(session, userID, exports.SecretGuildID) {
				// Remove the registered role and add the unregistered role
				err := session.GuildMemberRoleRemove(exports.SecretGuildID, userID, exports.RegisteredRole)
				if err != nil {
					log.Println("Error removing role from member:", err)
				}

				err = session.GuildMemberRoleAdd(exports.SecretGuildID, userID, exports.UnregisteredRole)
				if err != nil {
					log.Println("Error adding role to member:", err)
				}
			}

			_, err := session.ChannelMessageSend(exports.RegistrationChannel, fmt.Sprintf("<@%s> has been successfully unregistered.", userID)) // Added "successfully"
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
			_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> has been successfully unregistered.", userID)) // Added "successfully"
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		} else {
			_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> was not registered.", userID))
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}
	}
}
