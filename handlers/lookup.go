package handlers

import (
	"fmt"
	"strings"

	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func HandleWalletName(session *discordgo.Session, message *discordgo.MessageCreate) {
	loadUserMappings()
	content := message.Content
	// fmt.Println("CONTENT: %s", content)

	if !strings.HasPrefix(content, "!lookup ") {
		userMessage := "To lookup a DERO address, use the format: `!lookup <@user_mention or wallet_name>`"
		session.ChannelMessageSend(message.ChannelID, userMessage)
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!lookup ")
	// log.Printf("User Input: " + userInput)

	// Check if the input matches the format of a user mention
	if strings.HasPrefix(userInput, "<@") && strings.HasSuffix(userInput, ">") {
		userID := strings.TrimPrefix(userInput, "<@")
		userID = strings.TrimSuffix(userID, ">")

		userMappingsMutex.Lock()
		mappedAddress, exists := userMappings[userID]
		userMappingsMutex.Unlock()

		if exists {
			session.ChannelMessageSend(message.ChannelID, "DERO Address (from registered user): ```"+mappedAddress+"```")
		} else {
			// Ping the user with the mention
			userMention := "<@" + userID + ">"
			userMessage := " is not registered or invalid input. \n\n To register, please use `/register`"
			session.ChannelMessageSend(message.ChannelID, userMention+userMessage)
		}
	} else {
		// Perform a wallet name lookup for non-user mention inputs
		deroAddress, err := dero.WalletNameToAddress(userInput)
		if err != nil {
			fmt.Println("Wallet name not found or invalid.")
		}
		if deroAddress != "" {
			session.ChannelMessageSend(message.ChannelID, "DERO Address: ```"+deroAddress+"```")
		} else {
			session.ChannelMessageSend(message.ChannelID, "Wallet name not found or invalid.")
		}
	}
}
