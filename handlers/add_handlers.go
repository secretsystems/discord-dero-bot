package handlers

import (
	"fmt"
	"log"
	"strconv"

	"discord-dero-bot/utils/dero"

	"github.com/bwmarrin/discordgo"
)

var (
	ResultsChannel string
	session        string
)

func AddHandlers(discord *discordgo.Session, AppID, GuildID string) {
	// This handler will be triggered when the bot is ready
	log.Println("Registering Interaction Handlers")

	// Components are part of interactions, so we register InteractionCreate handler
	discord.AddHandler(func(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.Type {

		case discordgo.InteractionApplicationCommand:
			log.Println("received: discordgo.InteractionApplicationCommand")
			if h, ok := commandsHandlers[interaction.ApplicationCommandData().Name]; ok {
				h(discord, interaction, AppID, GuildID) // Pass appID and guildID
			}
		case discordgo.InteractionMessageComponent:
			log.Println("received: discordgo.InteractionMessageComponent")
			if h, ok := componentsHandlers[interaction.MessageComponentData().CustomID]; ok {
				h(discord, interaction, AppID, GuildID) // Pass appID and guildID
			}
		case discordgo.InteractionModalSubmit:
			data := interaction.ModalSubmitData()
			address := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			amountString := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			amount, error := strconv.Atoi(amountString)
			if error != nil {
				log.Printf("Error converting amount to int: %v", error)
			}
			comment := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			integratedAddress := dero.MakeIntegratedAddress(address, amount, comment)

			// Now you can use the integratedAddress
			fmt.Printf("Integrated Address: %s\n", integratedAddress)

			// Send an immediate response to the user
			err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: integratedAddress,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
			// not really sure why this doesn't work very well. there is some kind of panic:
			// HTTP 405 Method Not Allowed, {"message": "405: Method Not Allowed", "code": 0}
			// goroutine 40 [running]:
			// fuck_you.com/handlers.AddHandlers.func2(0xc00009f500?, 0xc000554020)
			//         /home/secret/phone_mount/code/fuck_you/handlers/add_handlers.go:67 +0x9ba
			// // Now you are going to be receiving this information some where
			// and so you are submiting it to something

			// if !strings.HasPrefix(data.CustomID, "encode") {
			// 	return
			// }
			// log.Printf("CustomID: %s", data.CustomID)

			// userid := strings.Split(data.CustomID, "_")[1]

			// // the results channel is defined in the config.go
			// _, err = discord.ChannelMessageSend(ResultsChannel, fmt.Sprintf("Request Received %v %v %v",
			// 	userid,

			// ),
			// )
			// if err != nil {
			// 	fmt.Println("I am panicing")
			// 	panic(err)
			// }

		}
	})

	// Register slash commands
	RegisterSlashCommands(discord, AppID, GuildID)
}
