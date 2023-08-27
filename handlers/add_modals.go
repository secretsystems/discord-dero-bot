package handlers

import (
	"discord-dero-bot/utils/dero"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var userMappings map[string]string
var userMappingsMutex sync.Mutex

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
	comment := ""
	comment += "C: " + color + " "
	comment += "S: " + size + " "
	comment += "A: " + shipping + " "
	comment += "P: " + contact + " "
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

	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for a Giftbox", userid)
	_, err = session.ChannelMessageSend(resultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}

func loadUserMappings() {
	file, err := os.OpenFile("userMappings.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening user mappings file: %v", err)
		return
	}
	defer file.Close()

	// Check if the file is empty before decoding
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return
	}
	if fileInfo.Size() == 0 {
		userMappings = make(map[string]string)
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&userMappings)
	if err != nil {
		log.Printf("Error decoding user mappings: %v", err)
	}
}

func saveUserMappings() {
	file, err := os.Create("userMappings.json")
	if err != nil {
		log.Printf("Error creating user mappings file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(userMappings)
	if err != nil {
		log.Printf("Error encoding user mappings: %v", err)
	}
}

func handleRegister(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

	// Respond to the interaction
	responseText := fmt.Sprintf("Successfully registered wallet input `%s` for <@%s>.", address, username)
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseText,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to interaction:", err)
	}
}
