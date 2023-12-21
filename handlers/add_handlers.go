package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	session string
)

func AddHandlers(session *discordgo.Session, appID string) {
	log.Println("Registering Interaction Handlers")

	handlers := DefineHandlers(session, appID)

	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.Type {
		case discordgo.InteractionApplicationCommand:
			log.Println("received: discordgo.InteractionApplicationCommand")
			if h, ok := handlers[interaction.ApplicationCommandData().Name]; ok {
				h(session, interaction, appID) // Pass appID
			}
		case discordgo.InteractionMessageComponent:
			log.Println("received: discordgo.InteractionMessageComponent")
			if h, ok := handlers[interaction.MessageComponentData().CustomID]; ok {
				h(session, interaction, appID) // Pass appID
			}
		}
	})
}
