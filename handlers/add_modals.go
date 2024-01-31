package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func AddModals(session *discordgo.Session, appID string) {
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Type == discordgo.InteractionModalSubmit {

			customID := interaction.ModalSubmitData().CustomID
			memberID := interaction.Member.User.ID

			// Distinguish between different custom IDs
			switch customID {
			case "encode_" + memberID:
				handleEncodeInteraction(session, interaction)
			case "decode_" + memberID:
				handleDecodeInteraction(session, interaction)
			case "giftbox_" + memberID:
				handleGiftboxInteraction(session, interaction)
			case "register_" + memberID:
				handleRegister(session, interaction)
			case "qr_" + memberID:
				handleQRInteraction(session, interaction)
			case "node_" + memberID:
				handleNodeInteraction(session, interaction)
			}
		}
	})
}
