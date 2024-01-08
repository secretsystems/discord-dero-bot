package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

var (
	owner              = "867976566716629092"
	registeredRoleID   = "1144842099653623839"
	unregisteredRoleID = "1144846590687838309"
	tipChannel         = "1161399751808385044"
	specialAddresses   = []string{
		"secret-wallet",
		"dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4",
	}
)

func init() {

	loadUserMap()

}

func HandleBigTip(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Create an array to store transfer information
	var transfers []dero.TransferInfo
	var discordIDs []string // To store Discord user IDs involved in transfers
	session.ChannelMessageSend(message.ChannelID, "All registered users are being tipped. This process takes time.")

	// Iterate through user mappings and create TransferInfo objects
	for discordID, address := range userMappings {
		address = resolveWalletAddress(address)
		amnt, _ := handleUserPermissions(session, message, discordID)
		transfer := dero.TransferInfo{
			Destination: address,
			Amount:      amnt, // Set the desired tip amount
		}
		transfers = append(transfers, transfer)

		// Add the Discord user ID to the list of users involved
		discordIDs = append(discordIDs, discordID)

		// If we have transfers, perform the bulk transfer and reset the transfers slice
		if len(transfers) == 10 {
			log.Printf("Before processing transfers: %v", transfers)
			processTransfers(session, message, transfers, discordIDs)
			log.Printf("After processing transfers")
			transfers = nil
			discordIDs = nil
		}
	}

	// Process any remaining transfers
	if len(transfers) > 0 {
		log.Printf("Before processing transfers: %v", transfers)
		processTransfers(session, message, transfers, discordIDs)
		log.Printf("After processing transfers")
	}

	return
}

func HandleTip(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content

	if content == "!tip" {
		session.ChannelMessageSend(message.ChannelID, "To send a tip, use the format: `!tip <@user_mention>`")
		return
	}

	if content == "!tips" {
		if message.Author.ID == owner {
			HandleBigTip(session, message)
		} else {
			session.ChannelMessageSend(message.ChannelID, "You don't have secret clearance.")
		}
	}

	if strings.HasPrefix(content, "!tip ") {
		input := strings.TrimPrefix(content, "!tip ")

		// Check for mentions and resolve user IDs
		mentionedUserIDs := resolveMentions(input)
		if len(mentionedUserIDs) > 0 {
			handleTip(session, message, mentionedUserIDs[1], getUserMappings(mentionedUserIDs[1]))
			return
		} else {
			handleTip(session, message, input, input)
			return
		}

	}
}

func processTransfers(session *discordgo.Session, message *discordgo.MessageCreate, transfers []dero.TransferInfo, discordIDs []string) {
	// Create a channel to communicate the result of the bulk transfer
	resultChan := make(chan string)
	defer close(resultChan) // Close the channel when done

	// Use a goroutine to make the bulk transfer concurrently
	go func() {
		// Call MakeBulkTransfer to perform the bulk transfer
		txID, err := dero.MakeBulkTransfer(transfers)
		if err != nil {
			resultChan <- "Error sending bulk tip: " + err.Error()
		} else {
			resultChan <- txID // Send the transaction ID to the channel
		}
	}()

	// Display the txid along with the success message and mention Discord users by their IDs
	messageToSend := "Bulk tip status:\n"
	txIDReceived := false

	// Wait for the result from the goroutine
	select {
	case result := <-resultChan:
		if strings.HasPrefix(result, "Error") {
			messageToSend += result
		} else {
			messageToSend += fmt.Sprintf("Explore this transaction by visiting http://explorer.dero.io/tx/%s"+
				"\nTips went to these registered badasses:\n", result)
			txIDReceived = true
		}
	}

	for _, discordID := range discordIDs {
		messageToSend += fmt.Sprintf("- <@%s>\n", discordID) // Mention users using Discord IDs
	}

	// Optionally, you can perform additional actions based on whether the bulk transfer was successful or not.
	if txIDReceived {
		// Send the message to the Discord channel
		session.ChannelMessageSend(tipChannel, messageToSend)
	} else {
		messageToSend = "TX FAILED TO BUILD TRANSFER"
		session.ChannelMessageSend(tipChannel, messageToSend) // Handle the case where the bulk transfer failed.
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

func handleMention(session *discordgo.Session, message *discordgo.MessageCreate, mentionedUserIDs []string) string {
	userID := mentionedUserIDs[1]
	mappedAddress := getUserMappings(userID)

	if mappedAddress == "" {
		// userMention := "<@" + userID + ">"
		// session.ChannelMessageSend(message.ChannelID, userMention+", you are not registered with tip bot.\n"+
		// 	"Please pair a DERO address with your profile by using the `/register` command.")
		return ""
	}

	return mappedAddress
}

func resolveWalletAddress(input string) string {
	if len(input) == 66 && strings.HasPrefix(input, "dero") {
		return input
	}

	lookupResult, err := dero.WalletNameToAddress(input) // Implement the wallet name lookup function
	if err != nil || lookupResult == "" {
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

func getAddressMappings(address string) string {
	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()
	return addressMappings[address]
}

func getUserMappings(userID string) string {
	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()
	return userMappings[userID]
}

func handleUserPermissions(session *discordgo.Session, message *discordgo.MessageCreate, userID string) (amnt int, amntmsg string) {

	amnt = 2
	amntmsg = "0.00002 DERO, or 2 DERI"
	// Get the user's roles
	member, err := session.GuildMember(secretGuildID, userID)
	if err != nil {
		log.Printf("Error getting guild member: %v", err)
		return
	}

	userRoles := member.Roles
	log.Printf("User has roles: %v", userRoles)

	// Check if the user ID is in userMappings
	if _, ok := userMappings[userID]; ok {
		// If the user ID is found in userMappings, adjust the tip amount accordingly
		amnt = 20 // Set the desired tip amount for registered users
		amntmsg = "0.00020 DERO, or 20 DERI"
		log.Printf("User ID: %s | Amount: %v | Message: %v", userID, amnt, amntmsg)
	}

	// Check user roles and adjust tip amount based on role priority
	for _, roleID := range userRoles {
		fmt.Printf("Role ID: %s\n", roleID)
		switch roleID {
		case secretMembersRoleID:
			amnt = 200
			amntmsg = "0.00200 DERO, or 200 DERI"
			log.Printf("Role ID: %s | Amount: %v | Message: %v", roleID, amnt, amntmsg)
		}
	}

	return amnt, amntmsg
}

func handleTip(session *discordgo.Session, message *discordgo.MessageCreate, userID, recipientAddress string) {
	var amnt int
	var amntmsg string

	switch {
	case userID == "":
		session.ChannelMessageSend(message.ChannelID, "userID cannot be empty")
		return
	case getUserMappings(userID) == "" && getAddressMappings(resolveWalletAddress(userID)) == "" && isValidDeroAddress(resolveWalletAddress(userID)) == false:
		session.ChannelMessageSend(message.ChannelID, "please consider using `/register`")
		return
	case getUserMappings(userID) != "":
		amnt, amntmsg = handleUserPermissions(session, message, userID)
		userID = fmt.Sprintf("<@%s>", userID)
	case getAddressMappings(resolveWalletAddress(userID)) != "":
		amnt, amntmsg = handleUserPermissions(
			session,
			message,
			getAddressMappings(
				resolveWalletAddress(userID),
			),
		)
		userID = fmt.Sprintf(
			"<@%s>",
			getAddressMappings(
				resolveWalletAddress(userID),
			),
		)

	case isValidDeroAddress(resolveWalletAddress(userID)) != false:
		amnt, amntmsg = handleUserPermissions(session, message, userID)
		userID = fmt.Sprintf("%s", userID)
	}

	waitingMessage := fmt.Sprintf("`secret-wallet` is sending %s to %s "+
		"This process takes roughly 18 seconds, or 1 block interval.",
		amntmsg,
		userID)

	session.ChannelMessageSend(message.ChannelID, waitingMessage)

	comment := "secret_pong_bot sends secret's love"

	txid, err := dero.MakeTransfer(resolveWalletAddress(recipientAddress), amnt, comment)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Error sending tip: "+err.Error())
		return
	}

	successMessage := fmt.Sprintf("TxID status for %s:\n```%s```\n"+
		"Explore this transaction by visiting: \n"+
		"> http://explorer.dero.io/tx/%s\n"+
		"Feed the bot by sending DERO to `secret-wallet`\n", userID, txid, txid)

	// Display the txid along with the success message
	session.ChannelMessageSend(tipChannel, successMessage)
}

func loadUserMap() error {
	data, err := os.ReadFile("userMappings.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &userMappings)
	if err != nil {
		return err
	}

	// Build the reverse map for address-based lookup
	for userID, address := range userMappings {
		addressMappings[address] = userID
	}

	// Ensure keys in the maps are trimmed and in lowercase
	// Modify the content of userMappings
	for k, v := range userMappings {
		delete(userMappings, k)
		userMappings[strings.TrimSpace(strings.ToLower(k))] = strings.TrimSpace(strings.ToLower(v))
	}

	return nil
}
