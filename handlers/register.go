package handlers

import (
	"discord-dero-bot/utils/dero"
	"fmt"
	"log"
	"strings"

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

	isDeroAddress := isValidDeroAddress(input)

	var address, walletType string

	if isDeroAddress {
		address = input
		walletType = "address"
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
		address = input
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
	if IsMemberInGuild(session, username, secretGuildID) {
		registeredRole := "1144842099653623839"
		err := session.GuildMemberRoleAdd(secretGuildID, username, registeredRole)
		if err != nil {
			log.Printf("Error adding role for Guild %v to member:%v", secretGuildID, err)
		}

		unregisteredRole := "1144846590687838309"
		err = session.GuildMemberRoleRemove(secretGuildID, username, unregisteredRole)
		if err != nil {
			log.Printf("Error removing role Guild %v to member:%v", secretGuildID, err)
		}
	}
	content := fmt.Sprintf("Successfully registered wallet %s `%s` for <@%s>.", walletType, address, username)
	RespondWithMessage(session, interaction, content)

	resultsChannel := "1156576030442655785"
	resultsMsg := fmt.Sprintf("<@%s> has registered with the server!", userID)
	_, err := session.ChannelMessageSend(resultsChannel, resultsMsg)
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

func isValidDeroAddress(address string) bool {
	return strings.HasPrefix(address, "dero1") && len(address) == 66
}
