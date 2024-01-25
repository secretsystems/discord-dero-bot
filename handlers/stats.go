package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/exports"
)

func HandleDerostats(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if the message starts with "!derostats"
	if !strings.HasPrefix(message.Content, "!derostats") {
		return
	}

	// Fetch JSON data from derostats.io
	data, err := getJSON()
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Failed to fetch JSON from derostats.io")
		return
	}
	response := "DEROSTATS:\n```"
	// Construct a single message containing all the information
	response += buildDerostatsInfo(data)
	response += "```\nsource: " + exports.DeroStatsURL
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

func getJSON() (map[string]interface{}, error) {

	// Send a GET request to the URL
	response, err := http.Get(exports.DeroStatsURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	// Decode the JSON response into a map
	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
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
