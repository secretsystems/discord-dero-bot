package handlers

import (
	"discord-dero-bot/utils/dero"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func AddModals(discord *discordgo.Session, AppID, GuildID string) {
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
			}
		}
	})

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
			Content: integratedAddress,
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
			Content: splitintegratedAddress,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Error responding to decode interaction:", err)
	}
}
