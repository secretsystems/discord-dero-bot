package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleShop(session *discordgo.Session, interaction *discordgo.MessageCreate) {
	helpCommand := strings.TrimPrefix(interaction.Content, "!shop ")

	switch helpCommand {
	case "list":
		// Send the general help message with a list of available commands
		helpMsg := "# Welcome to the Secret Discord Server Shop!\n"
		helpMsg += "## Available Product and Services(?):\n\n"
		helpMsg += "- !shop giftbox\n\n"

		session.ChannelMessageSend(interaction.ChannelID, helpMsg)
	case "giftbox":

		imagePath := "assets/giftbox.png"

		imageFile, err := os.Open(imagePath)
		if err != nil {
			log.Println("Error opening image path:", err)
			return
		}
		defer imageFile.Close()

		image := &discordgo.File{
			Name:        "giftbox.png",
			ContentType: "image/png",
			Reader:      imageFile,
		}
		content := "DERO GIFT BOX\n" +
			"Price: 55 USD w/ Shipping\n" +
			"The DERO gift box included:\n" +
			"- T-shirt with Dero logo\n" +
			"- 20 pens with Dero logo\n" +
			"- 4 dero collection stickers\n" +
			"	-  8 stickers with each collection\n" +
			"Delivered in 14-21 days\n" +
			"To order please, use ```/giftbox```"
		// Create a slice with the image file and assign it to the MessageSend struct
		messageSend := &discordgo.MessageSend{
			Files:   []*discordgo.File{image},
			Content: content,
		}

		_, err = session.ChannelMessageSendComplex(interaction.ChannelID, messageSend)
		if err != nil {
			log.Println("Error sending ephemeral message:", err)
		}
	// case "yogamat":

	// 	imagePath := "assets/giftbox.png"

	// 	imageFile, err := os.Open(imagePath)
	// 	if err != nil {
	// 		log.Println("Error opening image path:", err)
	// 		return
	// 	}
	// 	defer imageFile.Close()

	// 	image := &discordgo.File{
	// 		Name:        "giftbox.png",
	// 		ContentType: "image/png",
	// 		Reader:      imageFile,
	// 	}
	// 	content := "DERO GIFT BOX\n" +
	// 		"Price: 55 USD w/ Shipping\n" +
	// 		"The DERO gift box included:\n" +
	// 		"- T-shirt with Dero logo\n" +
	// 		"- 20 pens with Dero logo\n" +
	// 		"- 4 dero collection stickers\n" +
	// 		"	-  8 stickers with each collection\n" +
	// 		"Delivered in 14-21 days\n" +
	// 		"To order please, use ```/giftbox```"
	// 	// Create a slice with the image file and assign it to the MessageSend struct
	// 	messageSend := &discordgo.MessageSend{
	// 		Files:   []*discordgo.File{image},
	// 		Content: content,
	// 	}

	// 	_, err = session.ChannelMessageSendComplex(interaction.ChannelID, messageSend)
	// 	if err != nil {
	// 		log.Println("Error sending ephemeral message:", err)
	// 	}
	default:
		// Send an ephemeral message indicating the help command is not recognized
		helpMsg := "You have activated the !shop menu. Use `!shop list` to see available products."
		_, err := session.ChannelMessageSend(interaction.ChannelID, helpMsg)
		if err != nil {
			log.Println("Error sending ephemeral message:", err)
		}
	}
}
