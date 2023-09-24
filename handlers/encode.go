package handlers

import (
	"discord-dero-bot/utils/dero"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func handleEncodeModal(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
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

func handleEncodeInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ModalSubmitData()

	// Helper function to get a TextInput value by index
	getTextInputValue := func(index int) string {
		return data.Components[index].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	address := getTextInputValue(0)
	amountString := getTextInputValue(1)
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		log.Printf("Error converting amount to int: %v", err)
		RespondWithMessage(session, interaction, "Error: Invalid amount")
		return
	}
	comment := getTextInputValue(2)
	destinationString := getTextInputValue(3)
	destination, err := strconv.Atoi(destinationString)
	if err != nil {
		log.Printf("Error converting amount to int: %v", err)
		RespondWithMessage(session, interaction, "Error: Invalid destination")
		return
	}
	integratedAddress := dero.MakeIntegratedAddress(address, amount, comment, destination)
	RespondWithMessage(session, interaction, "```"+integratedAddress+"```")
}
