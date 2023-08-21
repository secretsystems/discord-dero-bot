package handlers

import (
	"log"
	"strings"

	"fuck_you.com/utils/dero"
	"github.com/bwmarrin/discordgo"
)

func HandleWalletName(discord *discordgo.Session, message *discordgo.MessageCreate) {
	userInput := strings.TrimPrefix(message.Content, "!lookup ")
	log.Printf("User Input: " + userInput)

	deroAddress := dero.WalletNameToAddress(userInput) // Use the new function to perform the lookup

	if deroAddress != "" {
		discord.ChannelMessageSend(message.ChannelID, "DERO Address: "+deroAddress)
	} else {
		discord.ChannelMessageSend(message.ChannelID, "Wallet name not found or invalid.")
	}
}
