package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/secretsystems/discord-dero-bot/exports"
	"github.com/secretsystems/discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

func handleRegistrationModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Command invoked in DM")
		RespondWithMessage(session, interaction, "This command cannot be used in DMs.")
		return
	}
	components := createRegisterModalComponents()
	modal := NewModal(session, interaction, "register_"+interaction.Interaction.Member.User.ID, "Secret Discord Server Registration", components)
	modal.Show()
}

func createRegisterModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "register",
				Label:       "Register DERO Address/Name with the Server",
				Style:       discordgo.TextInputShort,
				Placeholder: "dero1q wallet address or wallet-name",
				Required:    true,
			},
		}},
	}
}

func handleRegister(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Check if Member is nil (indicating DM)
	if interaction.Interaction.Member == nil {
		// Handle DM scenario
		log.Println("Interaction received in DM")
		RespondWithMessage(session, interaction, "This interaction cannot be processed in DMs.")
		return
	}
	loadUserMappings()
	data := interaction.ModalSubmitData()
	input := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	username := interaction.Member.User.ID
	isDeroAddress := isDeroAddress66Char(input)

	var address, walletType string

	if isDeroAddress {
		isValidDeroAddress := isValidDeroAddress(input)

		if isValidDeroAddress {
			address = input
			walletType = "address"
		} else {
			errMessage := "Error: wallet address not registered on DERO blockchain, please visit dero wallet and register address with DERO network"
			RespondWithMessage(session, interaction, errMessage)
			return
		}
	} else {
		walletType = "name"
		var err error
		address, err = dero.WalletNameToAddress(input)

		if err != nil {
			errMessage := fmt.Sprintf("Error: %s", err.Error())
			RespondWithMessage(session, interaction, errMessage)
			return
		}

		if address == "" {
			errMessage := fmt.Sprintf("<@%s> was looked up and is not registered to the DERO network. You need to register it with DERO first.", username)
			RespondWithMessage(session, interaction, errMessage)
			return
		}
	}

	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()

	if existingAddress, exists := userMappings[username]; exists {
		errMessage := fmt.Sprintf("<@%s> is already registered with wallet input: `%s`", username, existingAddress)
		RespondWithMessage(session, interaction, errMessage)
		return
	}

	for _, existingInput := range userMappings {
		if existingInput == address {
			errMessage := fmt.Sprintf("Wallet input `%s` is already registered.", address)
			RespondWithMessage(session, interaction, errMessage)

			return
		}
	}

	userMappings[username] = address
	saveUserMappings()

	// Get the user ID
	userID := strings.Split(data.CustomID, "_")[1]

	// Check if the interaction is in the desired guild (secretGuildID)
	if IsMemberInGuild(session, username, exports.SecretMembersRoleID) {

		err := session.GuildMemberRoleAdd(exports.SecretMembersRoleID, username, exports.RegisteredRole)
		if err != nil {
			log.Printf("Error adding role for Guild %v to member:%v", exports.SecretMembersRoleID, err)
		}

		err = session.GuildMemberRoleRemove(exports.SecretMembersRoleID, username, exports.UnregisteredRole)
		if err != nil {
			log.Printf("Error removing role Guild %v to member:%v", exports.SecretMembersRoleID, err)
		}
	}
	content := fmt.Sprintf("Successfully registered wallet %s `%s` for <@%s>.", walletType, address, username)
	RespondWithMessage(session, interaction, content)

	resultsMsg := fmt.Sprintf("<@%s> has registered with the server!", userID)
	_, err := session.ChannelMessageSend(exports.RegistrationChannel, resultsMsg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

func IsMemberInGuild(session *discordgo.Session, userID, guildID string) bool {
	member, err := session.GuildMember(guildID, userID)
	log.Println("Member is: ", member)
	if err != nil {
		log.Println("Error getting guild member:", err)
		return false
	}

	return member != nil
}

func isDeroAddress66Char(address string) bool {
	return strings.HasPrefix(address, "dero1") && len(address) == 66
}

func isValidDeroAddress(address string) bool {

	// Create a JSON request body
	requestBody := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": "1",
		"method": "DERO.GetEncryptedBalance",
		"params": {
			"address": "%s",
			"topoheight": -1
		}
	}`, address)

	// Define the DERO API URL
	apiURL := "http://" + dero.DeroServerIP + ":" + dero.DeroServerPort + "/json_rpc"

	// Make an HTTP POST request to the DERO API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		// Handle the error, e.g., log it
		log.Println("Error making HTTP request to DERO API:", err)
		return false
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Handle the non-200 response, e.g., log it
		log.Printf("DERO API returned non-200 status code: %d\n", resp.StatusCode)
		return false
	}

	// Parse the JSON response
	var response struct {
		JSONRPC string `json:"jsonrpc"`
		ID      string `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		// Handle the JSON parsing error, e.g., log it
		log.Println("Error decoding JSON response from DERO API:", err)
		return false
	}

	// Check if the address is valid based on the DERO API response
	if response.Error.Code == -32098 && response.Error.Message == "Account Unregistered" {
		// The address is unregistered
		return false
	} else if response.Error.Code == 0 && response.JSONRPC == "2.0" {
		// The address is registered
		return true
	}

	// Handle other unexpected cases
	log.Printf("Unexpected DERO API response: %+v\n", response)
	return false
}
