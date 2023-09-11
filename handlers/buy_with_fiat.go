package handlers

import (
	"discord-dero-bot/utils"

	"github.com/bwmarrin/discordgo"
)

func handleBuyDeroWithFiat(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
	// Create the content for the message
	content := createBuyDeroWithFiatContent()

	// Respond with the message containing trading information and payment link
	RespondWithMessage(session, interaction, content)
}

func createBuyDeroWithFiatContent() string {
	stripePaymentLink := "https://buy.stripe.com/7sI4k5bIp69HdI4cMN"
	return "DERO-USDT is trading at: " + utils.GetDeroUsdtAskString() + "\nWould you like to purchase DERO with Fiat?\nDisclosure: Limit of $100 and purchases have a fee of 8%\n" + stripePaymentLink
}
