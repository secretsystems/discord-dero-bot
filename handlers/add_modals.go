package handlers

import (
	"discord-dero-bot/utils"
	"discord-dero-bot/utils/dero"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddModals(session *discordgo.Session, appID, guildID, resultsChannel string) {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Check if the interaction type is ModalSubmitData
		if interaction.Type == discordgo.InteractionModalSubmit {
			// Get the CustomID from the interaction data
			customID := interaction.ModalSubmitData().CustomID

			// Distinguish between different custom IDs
			switch customID {
			case "encode_" + interaction.Member.User.ID:
				handleEncodeInteraction(session, interaction)
			case "decode_" + interaction.Member.User.ID:
				handleDecodeInteraction(session, interaction)
			case "giftbox_" + interaction.Member.User.ID:
				handleGiftboxInteraction(session, interaction, resultsChannel)
			case "register_" + interaction.Member.User.ID:
				handleRegister(session, interaction, resultsChannel)
			}
		}
	})
}

func handleEncodeInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ModalSubmitData()
	address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	amountString := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	amount, error := strconv.Atoi(amountString)
	if error != nil {
		log.Printf("Error converting amount to int: %v", error)
	}
	comment := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	destinationString := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	destination, error := strconv.Atoi(destinationString)
	if error != nil {
		log.Printf("Error converting amount to int: %v", error)
	}
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)

	// Now you can use the integratedAddress
	// fmt.Printf("Integrated Address: %s\n", integratedAddress)

	// Send an immediate response to the user

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "```" + integratedAddress + "```",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to encode interaction:", err)
	}
}

func handleDecodeInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ModalSubmitData()
	address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	splitintegratedAddress := dero.SplitIntegratedAddress(address)

	// Now you can use the integratedAddress
	// fmt.Printf("Integrated Address: %s\n", splitintegratedAddress)

	// Send an immediate response to the user
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "```" + splitintegratedAddress + "```",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to decode interaction:", err)
	}
}

func handleGiftboxInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate, resultsChannel string) {

	price := utils.ExchangeRate()
	// Step 3: Calculate the amount in atomic units
	amount := int((55 / price) * 100000)

	data := interaction.ModalSubmitData()

	color := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	size := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	shipping := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	contact := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	comment := ""
	comment += "C: " + color + " "
	comment += "S: " + size + " "
	comment += "A: " + shipping + " "
	comment += "P: " + contact + " "
	address := deroWalletAddress
	destination := 1337
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)
	messageContent := integratedAddress

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "To purchase your giftbox, please use the following address :\n```" + messageContent + "```And we will get back to you as soon as your order is marked receieved.\nWe will contact you on your shipping status.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Interaction Panic: %v", err)

		panic(err)
	}
	if !strings.HasPrefix(data.CustomID, "giftbox_") {
		return
	}

	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for a Giftbox", userid)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}

func handleRegister(session *discordgo.Session, interaction *discordgo.InteractionCreate, resultsChannel string) {
	loadUserMappings()
	data := interaction.ModalSubmitData()
	address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	username := interaction.Member.User.ID

	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()

	// Check if the user is already registered
	if existingAddress, exists := userMappings[username]; exists {
		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> is already registered with wallet input: `%s`", username, existingAddress),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Println("Error responding to decode interaction:", err)
		}
		return
	}

	// Check if the wallet input is already registered
	for _, existingInput := range userMappings {
		if existingInput == address {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Wallet input `%s` is already registered.", address),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("Error responding to decode interaction:", err)
			}
			return
		}
	}

	userMappings[username] = address
	saveUserMappings()

	// Add the registered role to the member
	err := session.GuildMemberRoleAdd(interaction.GuildID, username, registeredRole)
	if err != nil {
		log.Println("Error adding role to member:", err)
	}

	// Respond to the interaction
	responseText := fmt.Sprintf("Successfully registered wallet input `%s` for <@%s>.", address, username)
	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseText,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
	}

	userID := strings.Split(data.CustomID, "_")[1]

	resultsMsg := fmt.Sprintf("<@%s> has registered with the server!", userID)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		log.Println("Error sending message:", err) // Added log message
	}
}
