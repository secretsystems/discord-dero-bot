package handlers

// import (
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/bwmarrin/discordgo"
// )

// func HandleShop(session *discordgo.Session, interaction *discordgo.MessageCreate) {
// 	// Extract the command after "!shop"
// 	helpCommand := strings.TrimPrefix(interaction.Content, "!shop ")

// 	switch helpCommand {
// 	case "list":
// 		// Send a list of available products and services
// 		sendShopList(session, interaction.ChannelID)
// 	case "giftbox":
// 		// Send information about the giftbox and an image
// 		sendGiftboxInfo(session, interaction.ChannelID)
// 	default:
// 		// Send a message indicating that the command is not recognized
// 		sendShopHelp(session, interaction.ChannelID)
// 	}
// }

// func sendShopList(session *discordgo.Session, channelID string) {
// 	helpMsg := "# Welcome to the Secret Discord Server Shop!\n"
// 	helpMsg += "Available Product and Services(?):\n\n"
// 	helpMsg += "- !shop giftbox\n\n"
// 	_, err := session.ChannelMessageSend(channelID, helpMsg)
// 	if err != nil {
// 		log.Println("Error sending shop list:", err)
// 	}
// }

// func sendGiftboxInfo(session *discordgo.Session, channelID string) {
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
// 		"The DERO gift box includes:\n" +
// 		"- T-shirt with Dero logo\n" +
// 		"- 20 pens with Dero logo\n" +
// 		"- 4 DERO collection stickers\n" +
// 		"  - 8 stickers with each collection\n" +
// 		"Delivered in 14-21 days\n" +
// 		"To order, please use `/giftbox`"

// 	messageSend := &discordgo.MessageSend{
// 		Files:   []*discordgo.File{image},
// 		Content: content,
// 	}

// 	_, err = session.ChannelMessageSendComplex(channelID, messageSend)
// 	if err != nil {
// 		log.Println("Error sending giftbox info:", err)
// 	}
// }

// func sendShopHelp(session *discordgo.Session, channelID string) {
// 	helpMsg := "You have activated the !shop menu. Use `!shop list` to see available products."
// 	_, err := session.ChannelMessageSend(channelID, helpMsg)
// 	if err != nil {
// 		log.Println("Error sending shop help:", err)
// 	}
// }
