package handlers

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
)

// Handle !grok command
func HandleGrok(session *discordgo.Session, message *discordgo.MessageCreate) {
	imagePath := "assets/oogaboogaaa.jpg"
	imageFile, err := os.Open(imagePath)
	if err != nil {
		log.Println("Error opening image path:", err)
		return
	}
	defer imageFile.Close()

	image := &discordgo.File{
		Name:        "oogaboogaaa.jpg",
		ContentType: "image/png",
		Reader:      imageFile,
	}
	content := dero.GetGrok()

	messageSend := &discordgo.MessageSend{
		Files:   []*discordgo.File{image},
		Content: content,
	}

	_, err = session.ChannelMessageSendComplex(message.ChannelID, messageSend)
	if err != nil {
		log.Println("Error sending grokked info:", err)
	}
}
