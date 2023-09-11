package handlers

import (
	"discord-dero-bot/utils/dero"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleRegistration(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
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

	registeredRole := "1144842099653623839"
	err := session.GuildMemberRoleAdd(interaction.GuildID, username, registeredRole)
	if err != nil {
		log.Println("Error adding role to member:", err)
	}

	unregisteredRole := "1144846590687838309"
	err = session.GuildMemberRoleRemove(interaction.GuildID, username, unregisteredRole)
	if err != nil {
		log.Println("Error removing role from member:", err)
	}

	content := fmt.Sprintf("Successfully registered wallet %s `%s` for <@%s>.", walletType, address, username)
	RespondWithMessage(session, interaction, content)

	userID := strings.Split(data.CustomID, "_")[1]
	resultsChannel := "1060312629505167362"
	resultsMsg := fmt.Sprintf("<@%s> has registered with the server!", userID)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

func isValidDeroAddress(address string) bool {
	return strings.HasPrefix(address, "dero1") && len(address) == 66
}
