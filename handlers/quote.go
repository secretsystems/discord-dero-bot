package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleQuoteRequest(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	// fmt.Println("CONTENT: %s", content)
	if content == "!quote" {
		session.ChannelMessageSend(message.ChannelID, "To get a quote from TradeOgre: `!quote <base-pair>`\n\n For a list of pairs, use `!markets`")
		return

	} else if strings.HasPrefix(content, "!quote ") {
		ticker := strings.TrimPrefix(message.Content, "!quote ")
		// log.Printf("User Input: " + ticker)

		url := "https://tradeogre.com/api/v1/ticker/" + ticker

		// Create a GET request
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error sending HTTP Get request: %v", err)
			return
		}
		defer response.Body.Close()

		// Check if the response status code indicates an error
		if response.StatusCode != http.StatusOK {
			log.Printf("API request failed with status code: %d", response.StatusCode)
			session.ChannelMessageSend(message.ChannelID, "API request failed. Please try again later.")
			return
		}

		responseBody, _ := io.ReadAll(response.Body)
		// log.Printf("Response Body: %v", string(responseBody))

		var mapResponse map[string]interface{}
		err = json.Unmarshal(responseBody, &mapResponse)
		if err != nil {
			log.Printf("Error decoding response JSON: %v", err)
			return
		}

		// Print the entire mapResponse map
		// fmt.Println("\nmapResponse:", mapResponse)

		var outputMessage string
		for key, value := range mapResponse {
			formattedValue, _ := json.MarshalIndent(value, "", "  ")
			outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
		}

		// Send the entire response to session
		session.ChannelMessageSend(message.ChannelID, "Quote Response from TradeOgre.com:\n```\n"+outputMessage+"```")
	}
}
