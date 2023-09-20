package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHelp represents information about a specific command
type CommandHelp struct {
	Command  string
	Usage    string
	IsPublic bool
}

// HelpData stores information about available commands
var HelpData = []CommandHelp{
	{"!bot", "Get information from the bot based on your query", true},
	{"!unregister", "Unregister your wallet address or wallet name for tipping", true},
	{"!lookup", "Look up the DERO address associated with a username or wallet name", true},
	{"!tip", "Send a tip to the specified user or DERO address", true},
	{"!derod", "Get the current status of the DERO Network", true},
	{"!monerod", "Get the current status of the Monero Network", true},
	{"!markets", "Get the current list of markets provided by Trade Ogre", true},
	{"!quote", "Get the current quote of a pair provided by Trade Ogre", true},
	{"!derostats", "Get DERO stats from derostats.io", true},
	{"/trade-dero-xmr", "Trade DERO-XMR by way of the DERO-XMR swap integrated addresses", false},
	{"/encode", "Encode a DERO integrated Address using wallet address, amount, and comment", false},
	{"/decode", "Decode an integrated address and receive a DM of the output", false},
	{"/register", "Register your wallet address or wallet name for tipping", false},
}

func HandleHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	helpCommand := strings.TrimPrefix(message.Content, "!help ")

	if helpCommand == "list" {
		// Send a formatted list of available commands, separated by public (!) and private (/)
		helpMsg := "## Available Commands:\n"
		helpMsg += "### (!) Public Commands:\n```\n"
		for _, cmd := range HelpData {
			if cmd.IsPublic {
				helpMsg += fmt.Sprintf("%s: %s\n\n", cmd.Command, cmd.Usage)
			}
		}
		helpMsg += "```\n### (/) Private/Ephemeral Commands Please visit :\n```\n"
		for _, cmd := range HelpData {
			if !cmd.IsPublic {
				helpMsg += fmt.Sprintf("%s: %s\n\n", cmd.Command, cmd.Usage)
			}
		}
		helpMsg += "```\n"
		session.ChannelMessageSend(message.ChannelID, helpMsg)
		return
	}

	for _, cmd := range HelpData {
		if helpCommand == cmd.Command {
			// Send the usage information for the specific command and indicate public or private
			helpMsg := fmt.Sprintf("Usage: `%s`\n%s\n", cmd.Command, cmd.Usage)
			if cmd.IsPublic {
				helpMsg += "This is a public command."
			} else {
				helpMsg += "This is a private/ephemeral command."
			}
			session.ChannelMessageSend(message.ChannelID, helpMsg)
			return
		}
	}

	// If the specified command is not found, send a message indicating the help menu
	session.ChannelMessageSend(message.ChannelID, "You have activated the !help menu. Use `!help list` to see available commands.")
}
