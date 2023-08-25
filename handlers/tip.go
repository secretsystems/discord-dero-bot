package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"discord-dero-bot/utils/dero" // Import the dero package from your project

	"github.com/bwmarrin/discordgo"
)

func HandleTip(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	// fmt.Println("CONTENT: %s", content)
	if content == "!tip" {
		discord.ChannelMessageSend(message.ChannelID, "To send a tip, use the format: `!tip <@user_mention or <wallet_address> or <wallet_name>`")
		return

	} else if strings.HasPrefix(content, "!tip ") {

		// Extract the address or wallet name from the content
		input := strings.TrimPrefix(content, "!tip ")

		// Check if the input contains a mention
		mentionRegex := regexp.MustCompile("<@!?([0-9]+)>")
		mentionedUserIDs := mentionRegex.FindStringSubmatch(input)

		// Initialize variables to store user information
		var userID string
		var mappedAddress string
		var exists bool

		if len(mentionedUserIDs) == 2 {
			// A user was mentioned, look up their registered wallet address
			mentionedUserID := mentionedUserIDs[1]
			// fmt.Println("Mentioned User ID: %s", mentionedUserID)

			userMappingsMutex.Lock()
			mappedAddress, exists = userMappings[mentionedUserID]
			userMappingsMutex.Unlock()

			if exists {
				input = mappedAddress
			} else {
				userMention := "<@" + mentionedUserIDs[1] + ">"
				discord.ChannelMessageSend(message.ChannelID, userMention+", you are not registered with tip bot, please consider using `!register <wallet addr or wallet name>`")
				return
			}
		}

		// Check if the user input is in the userMappings
		userID = message.Author.ID
		userMappingsMutex.Lock()
		mappedAddress, exists = userMappings[userID]
		userMappingsMutex.Unlock()

		// Special addresses that should not receive tips
		specialAddresses := []string{
			"secret-wallet",
			"dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4",
		}

		// Check if the input address matches any special addresses
		for _, addr := range specialAddresses {
			if addr == input || (exists && addr == mappedAddress) {
				discord.ChannelMessageSend(message.ChannelID, "To tip the secret-bot, send funds to `secret-wallet`.")
				return
			}
		}

		// Determine the recipient's address for the tip
		var recipientAddress string

		// Check if the input is a valid DERO wallet address
		if len(input) == 66 && strings.HasPrefix(input, "dero") {
			recipientAddress = input
		} else {
			// Check if the input is a valid wallet name from the JSON
			if addr, ok := userMappings[input]; ok && len(addr) == 66 && strings.HasPrefix(addr, "dero") {
				recipientAddress = addr
			} else {
				// Perform a wallet name lookup
				lookupResult := dero.WalletNameToAddress(input) // Implement the wallet name lookup function

				if lookupResult != "" {
					// Ensure sender's address and recipient's address are different
					if lookupResult != mappedAddress {
						recipientAddress = lookupResult
					} else {
						discord.ChannelMessageSend(message.ChannelID, "To tip the secret-bot, send funds to `secret-wallet`.")
						return
					}
				} else {
					// Mention the mentioned user and provide the message
					if len(mentionedUserIDs) == 2 {
						userMention := "<@" + mentionedUserIDs[1] + ">"
						discord.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name.\n\n"+userMention+" Please consider using `!register <wallet addr or wallet name>`")
					} else {
						discord.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name.\n\nPlease consider using `!register <wallet addr or wallet name>`")
					}
					return
				}
			}
		}

		// Send the tip
		fmt.Println(recipientAddress)
		discord.ChannelMessageSend(message.ChannelID, "Tip is being sent from the `secret-wallet`: 0.00002 DERO, or 2 DERI")
		amnt := 2
		comment := "secret_pong_bot sends secret'a love"
		dero.MakeTransfer(recipientAddress, amnt, comment)
		discord.ChannelMessageSend(message.ChannelID, "Tip sent! \nFeed the bot by sending DERO to `secret-wallet`")
	}
}
