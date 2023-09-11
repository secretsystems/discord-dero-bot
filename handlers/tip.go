package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

var (
	secretMemberRoleID = "1057328486211145810"
	registeredRoleID   = "1144842099653623839"
	unregisteredRoleID = "1144846590687838309"
	specialAddresses   = []string{
		"secret-wallet",
		"dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4",
	}
)

func init() {
	loadUserMappings()
}

func HandleTip(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if content == "!tip" {
		session.ChannelMessageSend(message.ChannelID, "To send a tip, use the format: `!tip <@user_mention or <wallet_address> or <wallet_name>`")
		return
	}

	if strings.HasPrefix(content, "!tip ") {
		input := strings.TrimPrefix(content, "!tip ")

		// Check for mentions and resolve user IDs
		mentionedUserIDs := resolveMentions(input)
		if len(mentionedUserIDs) > 0 {
			handleMention(session, message, mentionedUserIDs)
			return
		}

		// Resolve wallet address or name
		recipientAddress := resolveWalletAddress(input, message.Author.ID)

		if recipientAddress == "" {
			session.ChannelMessageSend(message.ChannelID, "Invalid address or wallet name. Please consider using `/register`")
			return
		}

		// Check if the address is special
		if isSpecialAddress(recipientAddress) {
			session.ChannelMessageSend(message.ChannelID, "To tip the secret-bot, send funds to `secret-wallet`.")
			return
		}

		// Send the tip
		handleTip(session, message, recipientAddress)
	}
}

func resolveMentions(input string) []string {
	mentionRegex := regexp.MustCompile("<@!?([0-9]+)>")
	mentionedUserIDs := mentionRegex.FindStringSubmatch(input)
	if len(mentionedUserIDs) == 2 {
		return mentionedUserIDs
	}
	return nil
}

func handleMention(session *discordgo.Session, message *discordgo.MessageCreate, mentionedUserIDs []string) {
	userID := mentionedUserIDs[1]
	mappedAddress := getUserAddress(userID)

	if mappedAddress == "" {
		userMention := "<@" + userID + ">"
		session.ChannelMessageSend(message.ChannelID, userMention+", you are not registered with tip bot, please consider using `/register`")
		return
	}

	handleTip(session, message, mappedAddress)
}

func resolveWalletAddress(input, userID string) string {
	if len(input) == 66 && strings.HasPrefix(input, "dero") {
		return input
	}

	mappedAddress := getUserAddress(userID)
	if addr, ok := userMappings[input]; ok && len(addr) == 66 && strings.HasPrefix(addr, "dero") {
		return addr
	}

	lookupResult, err := dero.WalletNameToAddress(input) // Implement the wallet name lookup function
	if err != nil || lookupResult == "" || lookupResult == mappedAddress {
		return ""
	}

	return lookupResult
}

func isSpecialAddress(address string) bool {
	for _, addr := range specialAddresses {
		if addr == address {
			return true
		}
	}
	return false
}

func getUserAddress(userID string) string {
	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()
	return userMappings[userID]
}

func handleTip(session *discordgo.Session, message *discordgo.MessageCreate, recipientAddress string) {
	// Get the user's roles
	userRoles := message.Member.Roles

	// Define default tip amount and message
	amnt := 2
	amntmsg := "0.00002 DERO, or 2 DERI"

	// Check user roles and adjust tip amount based on role priority
	for _, roleID := range userRoles {
		switch roleID {
		case secretMemberRoleID:
			amnt = 200
			amntmsg = "0.00200 DERO, or 200 DERI"
		case registeredRoleID:
			if amnt != 200 {
				amnt = 20
				amntmsg = "0.00020 DERO, or 20 DERI"
			}
		case unregisteredRoleID:
			if amnt != 200 && amnt != 20 {
				amnt = 2
				amntmsg = "0.00002 DERO, or 2 DERI"
			}
		}
	}

	// Rest of the function remains the same as in the previous examples
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("`secret-wallet` is sending %s\nThis process takes roughly 18 seconds; or 1 block interval.", amntmsg))

	comment := "secret_pong_bot sends secret's love"
	txid, err := dero.MakeTransfer(recipientAddress, amnt, comment)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Error sending tip: "+err.Error())
		return
	}

	// Display the txid along with the success message
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Tip status:\n```TxID: %s```Feed the bot by sending DERO to `secret-wallet`", txid))
}
