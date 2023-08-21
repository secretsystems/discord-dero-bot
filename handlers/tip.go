package handlers

import (
	"regexp"
	"strings"

	"fuck_you.com/utils/dero" // Import the dero package from your project
	"github.com/bwmarrin/discordgo"
)

func HandleTip(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if strings.HasPrefix(content, "!tip ") {
		// Extract the address or wallet name from the content
		input := strings.TrimPrefix(content, "!tip ")

		// Check if the input contains a mention
		mentionRegex := regexp.MustCompile("<@!?([0-9]+)>")
		mentionedUserIDs := mentionRegex.FindStringSubmatch(input)

		// If a user was mentioned, look up their registered wallet address
		if len(mentionedUserIDs) == 2 {
			mentionedUserID := mentionedUserIDs[1]

			userMappingsMutex.Lock()
			mappedAddress, exists := userMappings[mentionedUserID]
			userMappingsMutex.Unlock()

			if exists {
				input = mappedAddress
			}
		}

		// Check if the user input is in the userMappings
		userID := message.Author.ID
		userMappingsMutex.Lock()
		mappedAddress, exists := userMappings[userID]
		userMappingsMutex.Unlock()

		// Special addresses that should not receive tips
		specialAddresses := []string{
			"secret-wallet",
			"dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4",
		}

		// Check if the input address matches any special addresses
		for _, addr := range specialAddresses {
			if addr == input || (exists && addr == mappedAddress) {
				discord.ChannelMessageSend(message.ChannelID, "You cannot send a tip to yourself.")
				return
			}
		}

		if exists && len(mappedAddress) == 66 && strings.HasPrefix(mappedAddress, "dero") {
			// If user is registered and their mapped address is valid, use it for the tip
			dero.MakeTransfer(mappedAddress)
			discord.ChannelMessageSend(message.ChannelID, "Tip sent!\n\n Please consider feeding the tip bot by sending DERO to `secret-wallet`")
		} else if len(input) == 66 && strings.HasPrefix(input, "dero") {
			// If input is a valid DERO address, use it directly for transfer
			dero.MakeTransfer(input)
			discord.ChannelMessageSend(message.ChannelID, "Tip sent!\n\n Please consider feeding the tip bot by sending DERO to `secret-wallet`")
		} else {
			// Otherwise, perform a wallet name lookup
			lookupResult := dero.WalletNameToAddress(input) // Implement the wallet name lookup function

			if lookupResult != "" {
				// Ensure sender's address and recipient's address are different
				if lookupResult != input {
					dero.MakeTransfer(lookupResult)
					discord.ChannelMessageSend(message.ChannelID, "Tip sent!\n\n Please consider feeding the tip bot by sending DERO to `secret-wallet`")
				} else {
					discord.ChannelMessageSend(message.ChannelID, "You cannot send a tip to yourself.")
				}
			} else {
				// Mention the mentioned user and provide the message
				if len(mentionedUserIDs) == 2 {
					userMention := "<@" + mentionedUserIDs[1] + ">"
					discord.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name.\n\n"+userMention+" Please consider using `!register`")
				} else {
					discord.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name.\n\nPlease consider using `!register`")
				}
			}
		}
	}
}
