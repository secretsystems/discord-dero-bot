package handlers

import (
	"log"
	"strings"

	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

var ()

func init() {

	loadUserMap()

}

func HandleVote(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	if content == "!vote" {
		if message.Author.ID == owner {
			session.ChannelMessageSend(message.ChannelID, "Please include scid.")
		} else {
			session.ChannelMessageSend(message.ChannelID, "You don't have secret clearance.")
		}
	}

	if strings.HasPrefix(content, "!vote ") {
		token := strings.TrimPrefix(content, "!vote ")

		// Check for mentions and resolve user IDs
		balance := hasBalance(token)
		if balance > 0 {
			distributeVoteTokens(session, message, token)
			return
		} else {
			session.ChannelMessageSend(message.ChannelID, "You don't have any vote tokens to distribute.")
			return
		}

	}
}

func distributeVoteTokens(session *discordgo.Session, message *discordgo.MessageCreate, token string) {
	// Create an array to store transfer information
	var transfers []dero.TransferInfo
	var discordIDs []string // To store Discord user IDs involved in transfers
	session.ChannelMessageSend(message.ChannelID, "All registered users are being tipped. This process takes time.")

	log.Printf("Before processing transfers: %v", transfers)
	// Iterate through user mappings and create TransferInfo objects
	for discordID, address := range userMappings {
		address = resolveWalletAddress(address)
		amnt, _ := handleUserPermissions(session, message, discordID)
		transfer := dero.TransferInfo{
			Destination: address,
			Amount:      amnt,
			SCID:        token,
		}
		transfers = append(transfers, transfer)

		// Add the Discord user ID to the list of users involved
		discordIDs = append(discordIDs, discordID)

		// If we have transfers, perform the bulk transfer and reset the transfers slice
		if len(transfers) == 50 {
			processTransfers(session, message, transfers, discordIDs)
			transfers = nil
			discordIDs = nil
		}
	}

	// Process any remaining transfers
	if len(transfers) > 0 {
		processTransfers(session, message, transfers, discordIDs)
	}
	log.Printf("After processing transfers")

	return
}

func hasBalance(scid string) int {
	return 0
}
