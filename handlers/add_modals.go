package handlers

import (
	"discord-dero-bot/utils/dero"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddModals(discord *discordgo.Session, AppID, GuildID string, ResultsChannel string) {
	discord.AddHandler(func(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Check if the interaction type is ModalSubmitData
		if interaction.Type == discordgo.InteractionModalSubmit {
			// Get the CustomID from the interaction data
			customID := interaction.ModalSubmitData().CustomID

			// Distinguish between different custom IDs
			switch customID {
			case "encode_" + interaction.Member.User.ID:
				handleEncodeInteraction(discord, interaction)
			case "decode_" + interaction.Member.User.ID:
				handleDecodeInteraction(discord, interaction)
			case "giftbox_" + interaction.Member.User.ID:
				handleGiftboxInteraction(discord, interaction, ResultsChannel)
			}
		}
	})
	log.Println("Adding Modals to Discord")
	// Register slash commands
	RegisterSlashCommands(discord, AppID, GuildID)
}

func handleEncodeInteraction(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

	err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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

func handleDecodeInteraction(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ModalSubmitData()
	address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	splitintegratedAddress := dero.SplitIntegratedAddress(address)

	// Now you can use the integratedAddress
	// fmt.Printf("Integrated Address: %s\n", splitintegratedAddress)

	// Send an immediate response to the user
	err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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

func handleGiftboxInteraction(discord *discordgo.Session, interaction *discordgo.InteractionCreate, ResultsChannel string) {
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

	// Step 3: Calculate the amount
	amount := int(55 / price)

	data := interaction.ModalSubmitData()

	color := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	size := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	shipping := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	contact := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	comment := ""
	comment += color
	comment += size
	comment += shipping
	comment += contact
	address := "dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4"
	destination := 1337
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)
	messageContent := integratedAddress

	err = discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "To purchase your giftbox, please use the following address :\n```" + messageContent + "```And we will get back to you as soon as your order is marked receieved.\nWe will contact you on your shipping status.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(data.CustomID, "giftbox_") {
		return
	}

	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for a Giftbox", userid)
	_, err = discord.ChannelMessageSend(ResultsChannel, resultsMsg)
	if err != nil {
		panic(err)
	}
}
