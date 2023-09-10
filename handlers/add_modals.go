package handlers

import (
	"discord-dero-bot/utils/coinbase"
	"discord-dero-bot/utils/dero"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddModals(session *discordgo.Session, appID, guildID string, resultsChannel string) {
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
				handleRegister(session, interaction)
			case "trade_dero_" + interaction.Member.User.ID:
				handleCryptoPurchase(session, interaction, resultsChannel)
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
	// Step 1: Make a GET request to the API endpoint
	response, err := http.Get("https://tradeogre.com/api/v1/ticker/dero-usdt")
	if err != nil {
		log.Println("Error fetching API data:", err)
		return
	}
	defer response.Body.Close()

	// Step 2: Parse the JSON response
	var apiResponse struct {
		Success bool   `json:"success"`
		Price   string `json:"price"`
	}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		log.Println("Error decoding API response:", err)
		return
	}

	// Convert the price from string to float64
	price, err := strconv.ParseFloat(apiResponse.Price, 64)
	if err != nil {
		log.Println("Error parsing price:", err)
		return
	}

	// Step 3: Calculate the amount in atomic units
	amount := int((55 / price) * 100000)

	data := interaction.ModalSubmitData()

	color := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	size := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	shipping := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	contact := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	comment := "C: " + color + " S: " + size + " A: " + shipping + " P: " + contact
	address := "dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4"
	destination := 1337
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)
	messageContent := integratedAddress

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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
	shopkeeper := "706842280828469260"
	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for <@%s>'s a Giftbox,", userid, shopkeeper)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}

func handleRegister(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	loadUserMappings()
	data := interaction.ModalSubmitData()
	input := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	username := interaction.Member.User.ID

	// Check if the input meets the criteria for a DERO address
	isDeroAddress := isValidDeroAddress(input)

	var address, walletType string

	if isDeroAddress {
		address = input
		walletType = "address"
	} else {
		walletType = "name"
		var err error
		address, err = dero.WalletNameToAddress(input) // Implement the wallet name lookup function

		if err != nil {
			errMessage := fmt.Sprintf("Error: %s", err.Error())
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: errMessage,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("Error responding to decode interaction:", err)
			}
			return
		}

		if address == "" {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> was looked up and is not registered to the DERO network. You need to register it with DERO first.", username),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("Error responding to decode interaction:", err)
			}
			return
		}
		address = input
	}

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

	registeredRole := "1144842099653623839"
	// Add the registered role to the member
	err := session.GuildMemberRoleAdd(interaction.GuildID, username, registeredRole)
	if err != nil {
		log.Println("Error adding role to member:", err)
	}

	unregisteredRole := "1144846590687838309"
	// Remove the unregistered role from the member
	err = session.GuildMemberRoleRemove(interaction.GuildID, username, unregisteredRole)
	if err != nil {
		log.Println("Error removing role from member:", err) // Updated log message
	}

	// Respond to the interaction
	responseText := fmt.Sprintf("Successfully registered wallet %s `%s` for <@%s>.", walletType, address, username)
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
	resultsChannel := "1060312629505167362"
	resultsMsg := fmt.Sprintf("<@%s> has registered with the server!", userID)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		log.Println("Error sending message:", err) // Added log message
	}
}

func handleCryptoPurchase(session *discordgo.Session, interaction *discordgo.InteractionCreate, resultsChannel string) {
	// Step 1: Make a GET request to the API endpoint
	resultsChannel = "1059682504124158074"
	response, err := http.Get("https://tradeogre.com/api/v1/ticker/dero-usdt")
	if err != nil {
		log.Println("Error fetching API data:", err)
		return
	}
	defer response.Body.Close()

	// Step 2: Parse the JSON response
	var apiResponse struct {
		Success bool   `json:"success"`
		Price   string `json:"price"`
	}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		log.Println("Error decoding API response:", err)
		return
	}

	// Convert the price from string to float64
	price, err := strconv.ParseFloat(apiResponse.Price, 64)
	if err != nil {
		log.Println("Error parsing price:", err)
		return
	}

	// Step 3: Calculate the amount in atomic units

	data := interaction.ModalSubmitData()

	amountString := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		log.Printf("Error parsing price: %s", err)
		return
	}
	price = price * amount
	address := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	messageContent := coinbase.PostCharges(price)

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Please visit the coinbase address to complete your purchase :\n " + messageContent + " \nAnd we will get back to you as soon as your order is marked receieved.\nWe will contact you on your shipping status.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Interaction Panic: %v", err)

		panic(err)
	}
	if !strings.HasPrefix(data.CustomID, "trade_dero_") {
		return
	}

	userid := strings.Split(data.CustomID, "_")[2]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an order with address: %s", userid, address)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}
func isValidDeroAddress(address string) bool {
	// Check if the address starts with "dero1" and has 66 characters
	return strings.HasPrefix(address, "dero1") && len(address) == 66
}
