package handlers

import (
	"log"

	"fuck_you.com/commands"
	"github.com/bwmarrin/discordgo"
)

func RegisterInteractionHandlersFromHandlers(discord *discordgo.Session, AppID, GuildID string) {
	// This handler will be triggered when the bot is ready
	log.Println("Registering Interaction Handlers")
	discord.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	// Components are part of interactions, so we register InteractionCreate handler
	discord.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.Type {
		case discordgo.InteractionApplicationCommand:
			log.Println("recieved: discordgo.InteractionApplicationCommand")
			if h, ok := commandsHandlers[interaction.ApplicationCommandData().Name]; ok {
				h(session, interaction, AppID, GuildID) // Pass appID and guildID
			}
		case discordgo.InteractionMessageComponent:
			log.Println("recieved: discordgo.InteractionMessageComponent")
			if h, ok := componentsHandlers[interaction.MessageComponentData().CustomID]; ok {
				h(session, interaction, AppID, GuildID) // Pass appID and guildID
			}
		}
	})

	// Register slash commands
	commands.RegisterSlashCommands(discord, AppID, GuildID)
}
