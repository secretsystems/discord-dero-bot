package handlers

import "github.com/bwmarrin/discordgo"

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
