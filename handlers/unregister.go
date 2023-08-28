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
		userID := message.Author.ID
		log.Printf("User ID: %s", userID) // Added "User ID:"

		userMappingsMutex.Lock()
		defer userMappingsMutex.Unlock()

		// Check if the user is registered
		if _, exists := userMappings[userID]; exists {
			delete(userMappings, userID)
			saveUserMappings()

			registeredRole := "1144842099653623839"
			unregisteredRole := "1144846590687838309"

			// Remove the registered role and add the unregistered role
			err := session.GuildMemberRoleRemove(message.GuildID, userID, registeredRole)
			if err != nil {
				log.Println("Error removing role from member:", err) // Updated log message
			}

			err = session.GuildMemberRoleAdd(message.GuildID, userID, unregisteredRole)
			if err != nil {
				log.Println("Error adding role to member:", err)
			}

			_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> has been successfully unregistered.", userID)) // Added "successfully"
			if err != nil {
				log.Printf("Error sending message: %v\n", err) // Used log.Printf
			}
		} else {
			_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> was not registered.", userID))
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}
	}
}
