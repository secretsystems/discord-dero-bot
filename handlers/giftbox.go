package handlers

import "github.com/bwmarrin/discordgo"

func handleGiftbox(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
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
