package handlers

import (
	"strings"

	"discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func HandleWalletName(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	// fmt.Println("CONTENT: %s", content)

	if !strings.HasPrefix(content, "!lookup ") {
		discord.ChannelMessageSend(message.ChannelID, "To lookup a DERO address, use the format: `!lookup <@user_mention or wallet_name>`")
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
			discord.ChannelMessageSend(message.ChannelID, "DERO Address (from registered user): ```"+mappedAddress+"```")
		} else {
			// Ping the user with the mention
			userMention := "<@" + userID + ">"
			discord.ChannelMessageSend(message.ChannelID, userMention+" is not registered or invalid input. \n\n To register, please use `!register`")
		}
	} else {
		// Perform a wallet name lookup for non-user mention inputs
		deroAddress := dero.WalletNameToAddress(userInput)

		if deroAddress != "" {
			discord.ChannelMessageSend(message.ChannelID, "DERO Address: ```"+deroAddress+"```")
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Wallet name not found or invalid.")
		}
	}
}
