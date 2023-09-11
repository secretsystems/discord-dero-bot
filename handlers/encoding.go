package handlers

import "github.com/bwmarrin/discordgo"

func handleEncode(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID, guildID string) {
	components := createEncodeModalComponents()
	modal := NewModal(session, interaction, "encode_"+interaction.Interaction.Member.User.ID, "Encode a DERO Integrated Address", components)
	modal.Show()
}
func createEncodeModalComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "address",
				Label:       "Address of where funds will be sent",
				Style:       discordgo.TextInputShort,
				Placeholder: "dero1q wallet address",
				Required:    true,
				MaxLength:   66,
				MinLength:   66,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "amount",
				Label:       "Amount in atomic units; minimum 2 DERI",
				Style:       discordgo.TextInputShort,
				Placeholder: "1 DERO = 100000 ; 2 DERI = 2",
				Required:    true,
				MaxLength:   64,
				MinLength:   1,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "comment",
				Label:       "Comment/subject/details",
				Style:       discordgo.TextInputParagraph,
				Placeholder: "",
				Required:    false,
				MaxLength:   128,
			},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "destination",
				Label:       "What port you want to send this to?",
				Style:       discordgo.TextInputShort,
				Placeholder: "ex. 1337",
				Required:    false,
				MaxLength:   128,
			},
		}},
	}
}
