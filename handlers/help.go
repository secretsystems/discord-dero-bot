package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	helpCommand := strings.TrimPrefix(message.Content, "!help ")

	switch helpCommand {
	case "list":
		// Send the general help message with a list of available commands
		helpMsg := "# Welcome to the Secret Discord Server!\n"
		helpMsg += "```We share our knowledge, insights and relationships we earned from our research and development using DERO.```\n"
		helpMsg += "## Available Commands:\n\n"
		helpMsg += "**Server:**\n"
		helpMsg += "```!help <!command>``````!bot <query>``````/register``````!unregister``````/trade-dero-xmr```\n"
		helpMsg += "**DERO Wallet:**\n"
		helpMsg += "```!lookup <@username or wallet name>``````!tip <@username, dero1q or wallet-name>``````/decode``````/encode```\n"
		helpMsg += "**Node**\n"
		helpMsg += "```!derod``````!monerod```\n"
		helpMsg += "**Markets**\n"
		helpMsg += "```!markets``````!quote <insert base-pair>```\n"
		session.ChannelMessageSend(message.ChannelID, helpMsg)
	case "!bot":
		// Send a breakdown of the bot command and its usage
		botHelpMsg := "Usage: `!bot <query>`\n" +
			"Get information from the bot based on your query."
		session.ChannelMessageSend(message.ChannelID, botHelpMsg)
	case "!register":
		// Send a breakdown of the register command and its usage
		registerHelpMsg := "Usage: `!register <wallet address or wallet-name>`\n" +
			"Register your wallet address or wallet name for tipping."
		session.ChannelMessageSend(message.ChannelID, registerHelpMsg)
	case "!unregister":
		// Send a breakdown of the register command and its usage
		unregisterHelpMsg := "Usage: `!unregister`\n" +
			"Unregister your wallet address or wallet name for tipping."
		session.ChannelMessageSend(message.ChannelID, unregisterHelpMsg)
	case "!lookup":
		// Send a breakdown of the lookup command and its usage
		lookupHelpMsg := "Usage: `!lookup <@username or wallet name>`\n" +
			"Look up the DERO address associated with a username or wallet name."
		session.ChannelMessageSend(message.ChannelID, lookupHelpMsg)
	case "/decode":
		// Send a breakdown of the lookup command and its usage
		lookupHelpMsg := "Usage: `/decode`\n" +
			"Decode an integrated address and receive a DM of the output."
		session.ChannelMessageSend(message.ChannelID, lookupHelpMsg)
	case "!tip":
		// Send a breakdown of the tip command and its usage
		tipHelpMsg := "Usage: `!tip <@username, dero1q, or wallet-name>`\n" +
			"Send a tip to the specified user or DERO address."
		session.ChannelMessageSend(message.ChannelID, tipHelpMsg)
	case "!derod":
		// Send a breakdown of the derod command and its usage
		derodHelpMsg := "Usage: `!derod`\n" +
			"Get the current status of the DERO Network."
		session.ChannelMessageSend(message.ChannelID, derodHelpMsg)
	case "!monerod":
		// Send a breakdown of the monerod command and its usage
		monerodHelpMsg := "Usage: `!monerod`\n" +
			"Get the current status of the Monero Network."
		session.ChannelMessageSend(message.ChannelID, monerodHelpMsg)
	case "!markets":
		// Send a breakdown of the markets command and its usage
		marketsHelpMsg := "Usage: !markets\n" +
			"Get the current list of markets provided by Trade Ogre."
		session.ChannelMessageSend(message.ChannelID, marketsHelpMsg)
	case "!quote":
		// Send a breakdown of the quote command and its usage
		quoteHelpMsg := "Usage: `!quote <insert base-pair>`\n" +
			"Get the current quote of a pair provided by Trade Ogre."
		session.ChannelMessageSend(message.ChannelID, quoteHelpMsg)
	case "/trade-dero-xmr":
		// Send a breakdown of the quote command and its usage
		tradeHelpMsg := "Usage: `/trade-dero-xmr`\n" +
			"Trade DERO-XMR by way of the DERO-XMR swap integrated addresses."
		session.ChannelMessageSend(message.ChannelID, tradeHelpMsg)
	case "/encode":
		// Send a breakdown of the quote command and its usage
		tradeHelpMsg := "Usage: `/encode`\n" +
			"Encode a DERO integrated Address using wallet address, amount and comment."
		session.ChannelMessageSend(message.ChannelID, tradeHelpMsg)
	default:
		// Send a message indicating the help command is not recognized
		session.ChannelMessageSend(message.ChannelID, "You have activated the !help menu. Use `!help list` to see available commands.")
	}
}
