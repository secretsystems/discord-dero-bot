package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleHelp(discord *discordgo.Session, message *discordgo.MessageCreate) {
	helpCommand := strings.TrimPrefix(message.Content, "!help ")

	switch helpCommand {
	case "!me":
		// Send the general help message with a list of available commands
		helpMsg := "# Welcome to the Secret Discord Server!\n"
		helpMsg += "```We share our knowledge, insights and relationships we earned from our research and development using DERO.```\n"
		helpMsg += "## Available Commands:\n\n"
		helpMsg += "### `!commands` List:\n\n"
		helpMsg += "**Server:**\n"
		helpMsg += "```!help <!command>``````!bot <query>``````!register <wallet address or wallet-name>``````!unregister```\n"
		helpMsg += "**DERO Wallet:**\n"
		helpMsg += "```!decode <integrated address>``````!lookup <@username or wallet name>``````!tip <@username, dero1q or wallet-name>```\n"
		helpMsg += "**Node**\n"
		helpMsg += "```!derod``````!monerod```\n"
		helpMsg += "**Markets**\n"
		helpMsg += "```!markets``````!quote <insert base-pair>```\n"
		helpMsg += "### `/commands` List:\n"
		helpMsg += "\n"
		helpMsg += "**Trades**\n"
		helpMsg += "```/trade-dero-xmr```\n"
		discord.ChannelMessageSend(message.ChannelID, helpMsg)
	case "!bot":
		// Send a breakdown of the bot command and its usage
		botHelpMsg := "Usage: !bot <query>\n" +
			"Get information from the bot based on your query."
		discord.ChannelMessageSend(message.ChannelID, botHelpMsg)
	case "!register":
		// Send a breakdown of the register command and its usage
		registerHelpMsg := "Usage: !register <wallet address or wallet-name>\n" +
			"Register your wallet address or wallet name for tipping."
		discord.ChannelMessageSend(message.ChannelID, registerHelpMsg)
	case "!lookup":
		// Send a breakdown of the lookup command and its usage
		lookupHelpMsg := "Usage: !lookup <@username or wallet name>\n" +
			"Look up the DERO address associated with a username or wallet name."
		discord.ChannelMessageSend(message.ChannelID, lookupHelpMsg)
	case "!tip":
		// Send a breakdown of the tip command and its usage
		tipHelpMsg := "Usage: !tip <@username, dero1q, or wallet-name>\n" +
			"Send a tip to the specified user or DERO address."
		discord.ChannelMessageSend(message.ChannelID, tipHelpMsg)
	case "!derod":
		// Send a breakdown of the derod command and its usage
		derodHelpMsg := "Usage: !derod\n" +
			"Get the current status of the DERO Network."
		discord.ChannelMessageSend(message.ChannelID, derodHelpMsg)
	case "!monerod":
		// Send a breakdown of the monerod command and its usage
		monerodHelpMsg := "Usage: !monerod\n" +
			"Get the current status of the Monero Network."
		discord.ChannelMessageSend(message.ChannelID, monerodHelpMsg)
	case "!markets":
		// Send a breakdown of the markets command and its usage
		marketsHelpMsg := "Usage: !markets\n" +
			"Get the current list of markets provided by Trade Ogre."
		discord.ChannelMessageSend(message.ChannelID, marketsHelpMsg)
	case "!quote":
		// Send a breakdown of the quote command and its usage
		quoteHelpMsg := "Usage: !quote <insert base-pair>\n" +
			"Get the current quote of a pair provided by Trade Ogre."
		discord.ChannelMessageSend(message.ChannelID, quoteHelpMsg)
	default:
		// Send a message indicating the help command is not recognized
		discord.ChannelMessageSend(message.ChannelID, "You have activated the !help menu. Use `!help !me` to see available commands.")
	}
}
