package handlers

import (
	"fmt"
	"github.com/secretsystems/discord-dero-bot/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleDerostats(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if the message starts with "!derostats"
	if !strings.HasPrefix(message.Content, "!derostats") {
		return
	}

	// Fetch JSON data from derostats.io
	data, err := utils.GetJSON()
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Failed to fetch JSON from derostats.io")
		return
	}
	response := "DEROSTATS:\n```"
	// Construct a single message containing all the information
	response += buildDerostatsInfo(data)
	response += "```\nsource: https://derostats.io"
	session.ChannelMessageSend(message.ChannelID, response)
}

func buildDerostatsInfo(data map[string]interface{}) string {
	var response strings.Builder

	appendSection := func(title string, sectionData interface{}) {
		response.WriteString(title)
		response.WriteString(":\n")

		switch sectionData := sectionData.(type) {
		case map[string]interface{}:
			for key, value := range sectionData {
				if key == "updated" || key == "block" {
					// Convert large numbers to a more human-readable format
					convertedValue := convertNumberToReadable(value)
					response.WriteString(fmt.Sprintf("%s: %s\n", key, convertedValue))
				} else {
					response.WriteString(fmt.Sprintf("%s: %v\n", key, value))
				}
			}
		case float64:
			response.WriteString(fmt.Sprintf("%.0f\n", sectionData))
		}

		response.WriteString("\n")
	}

	appendSection("Information", data["informations"])
	appendSection("Fees", data["fees"])
	appendSection("Total Accounts", data["accounts"])
	appendSection("Total Transactions", data["transactions"])
	appendSection("Smart Contract Executions", data["scExecutions"])
	appendSection("Smart Contracts Deployed", data["scDeployed"])

	return response.String()
}

func convertNumberToReadable(value interface{}) string {
	// Convert the large numbers to a more human-readable format
	switch value := value.(type) {
	case float64:
		return fmt.Sprintf("%.0f", value)
	case string:
		return value
	default:
		return fmt.Sprintf("%v", value)
	}
}
