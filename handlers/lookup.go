package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleWalletName(session *discordgo.Session, message *discordgo.MessageCreate) {
	loadUserMappings()
	content := message.Content
	// fmt.Println("CONTENT: %s", content)

	userInput := strings.TrimPrefix(content, "!lookup ")
	// log.Printf("User Input: " + userInput)

	switch {
	case userInput == "":
		userMessage := "To lookup a DERO address, use the format: `!lookup <@user_mention or wallet_name>`"
		session.ChannelMessageSend(message.ChannelID, userMessage)
		return
	case strings.HasPrefix(userInput, "<@"):
		userID := strings.TrimPrefix(userInput, "<@")
		userID = strings.TrimSuffix(userID, ">")

		exists := getUserMappings(userID)

		if exists != "" {
			session.ChannelMessageSend(message.ChannelID, "DERO Address: ```"+exists+"```")
		}
	case getUserMappings(userInput) == "" && getAddressMappings(resolveWalletAddress(userInput)) == "" && isValidDeroAddress(resolveWalletAddress(userInput)):
		session.ChannelMessageSend(message.ChannelID, "please consider using `/register`")
		return
	case getAddressMappings(resolveWalletAddress(userInput)) != "":
		userInput = fmt.Sprintf(
			"<@%s>",
			getAddressMappings(
				resolveWalletAddress(userInput),
			),
		)
		session.ChannelMessageSend(message.ChannelID, "Discord User: "+userInput)

	}
}
