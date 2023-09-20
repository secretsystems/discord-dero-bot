package handlers

import (
	"discord-dero-bot/utils"
	"discord-dero-bot/utils/dero"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleGiftbox(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	components := createGiftBoxModalComponents()
	modal := NewModal(session, interaction, "giftbox_"+interaction.Interaction.Member.User.ID, "Purchase a DERO Gift Box", components)
	modal.Show()
}
func createGiftBoxModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "color",
				Label:       "Shirt Color",
				Style:       discordgo.TextInputShort,
				Placeholder: "black or white?",
				Required:    true,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "size",
				Label:       "Shirt Size",
				Style:       discordgo.TextInputShort,
				Placeholder: "What size fits you: S, M, L, XL, XXL, XXXL",
				Required:    true,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:  "address",
				Label:     "Shipping Address?",
				Style:     discordgo.TextInputParagraph,
				Required:  false,
				MaxLength: 80,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:  "contact_info",
				Label:     "If we need to reach you?",
				Style:     discordgo.TextInputShort,
				Required:  false,
				MaxLength: 128,
			},
		}},
	}
}

func handleGiftboxInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	price := utils.GetAsk("dero-usdt")
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

	content := "To purchase your giftbox, please use the following address :\n```" + integratedAddress + "```And we will get back to you as soon as your order is marked receieved.\nWe will contact you on your shipping status."
	RespondWithMessage(session, interaction, content)

	if !strings.HasPrefix(data.CustomID, "giftbox_") {
		return
	}
	shopkeeper := "706842280828469260"
	userid := strings.Split(data.CustomID, "_")[1]
	resultsMsg := fmt.Sprintf(
		"User <@%s> has made an integrated address for <@%s>'s a Giftbox,", userid, shopkeeper)
	_, err := session.ChannelMessageSend(interaction.ChannelID, resultsMsg)
	if err != nil {
		panic(err)
	}
}
