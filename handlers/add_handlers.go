package handlers

import (
	"log"

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
		}
	})
	log.Println("Adding Handlers to Discord")
	// Register slash commands
	RegisterSlashCommands(discord, AppID, GuildID)
}
