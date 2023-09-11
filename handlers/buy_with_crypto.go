package handlers

import (
	"discord-dero-bot/utils"
	"discord-dero-bot/utils/coinbase"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleBuyDeroWithCrypto(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
	// Define the components for the Purchase with Crypto modal
	components := createBuyDeroWithCryptoModalComponents()

	// Respond with the Purchase with Crypto modal
	modal := NewModal(session, interaction, "trade_dero_"+interaction.Interaction.Member.User.ID, "Purchase DERO with Crypto", components)
	modal.Show()
}

func createBuyDeroWithCryptoModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "amount",
				Label:       "How much DERO would you like",
				Style:       discordgo.TextInputShort,
				Placeholder: "Please keep in amounts under 100 DERO",
				Required:    true,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "address",
				Label:       "Wallet Address",
				Style:       discordgo.TextInputShort,
				Placeholder: "Please provide full dero1q address",
				Required:    false,
				MaxLength:   66,
				MinLength:   66,
			},
		}},
	}
}

func handleCryptoPurchase(session *discordgo.Session, interaction *discordgo.InteractionCreate, resultsChannel string) {
	resultsChannel = "1059682504124158074"

	price := utils.GetAsk("dero-usdt")

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

	content := "Please visit the coinbase address to complete your purchase :\n " + messageContent + " \nAnd we will get back to you as soon as your order is marked receieved.\nWe will contact you on your status."
	RespondWithMessage(session, interaction, content)
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
