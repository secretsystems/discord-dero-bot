package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func AddModals(session *discordgo.Session, appID, guildID string, resultsChannel string) {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Type == discordgo.InteractionModalSubmit {

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
