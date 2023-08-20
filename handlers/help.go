package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleHelp(discord *discordgo.Session, message *discordgo.MessageCreate) {
	helpMessage := "Welcome to the Secret Discord Server!\n"
	helpMessage += "```We share our knowledge, insights and relationships we earned from our research and development using DERO.```\n"
	helpMessage += "Available Commands:\n\n"
	helpMessage += "`!commands` List:\n\n"
	helpMessage += "**For privacy, you can dm the bot.**\n\n"
	helpMessage += "Server:\n"
	helpMessage += "```!help```"
	// helpMessage += "Produces this help message\n"
	helpMessage += "\n"
	helpMessage += "DERO Wallet:\n"
	helpMessage += "```!decode <insert integrated address>```"
	// helpMessage += "You can decode a integrated address\n"
	helpMessage += "```!lookup <insert wallet name>```"
	// helpMessage += "You can look up a wallet name address\n"
	helpMessage += "\n"
	helpMessage += "Node\n"
	helpMessage += "```!derod```"
	// helpMessage += "Get the current status of the DERO Network \n"
	helpMessage += "```!monerod```"
	// helpMessage += "Get the current status of the Monero Network \n"
	helpMessage += "\n"
	helpMessage += "Markets\n"
	helpMessage += "```!markets```"
	// helpMessage += "Get the current list of markets provided by Trade Ogre\n"
	helpMessage += "```!quote <insert base-pair>```"
	// helpMessage += "Get the current quote of a pair provided by Trade Ogre\n"
	helpMessage += "\n"
	helpMessage += "`/commands` List:\n"
	helpMessage += "\n"
	helpMessage += "Trades\n"
	helpMessage += "```/trade-dero-xmr```"
	// helpMessage += "You can decode a integrated address\n"

	discord.ChannelMessageSend(message.ChannelID, helpMessage)
}
