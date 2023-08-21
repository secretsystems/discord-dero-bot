// handlers/integrated_address.go

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleIntegratedAddress(discord *discordgo.Session, message *discordgo.MessageCreate) {
	userInput := strings.TrimPrefix(message.Content, "!decode ")
	log.Printf("User Input: " + userInput)

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "SplitIntegratedAddress",
		"params": map[string]string{
			"integrated_address": userInput,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return
	}

	request, err := http.NewRequest("POST", "http://192.168.12.208:10103/json_rpc", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}

	// Set Basic Authentication
	request.SetBasicAuth("user", "pass")
	request.Header.Set("Content-type", "application/json")
	fmt.Println("\nRequest: ", request)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	var mapResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &mapResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return
	}

	// Print the entire mapResponse map
	fmt.Println("\nmapResponse:", mapResponse)

	// Print the entire mapResponse map
	var outputMessage string
	for key, value := range mapResponse {
		formattedValue, _ := json.MarshalIndent(value, "", "  ")
		outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
	}

	// Send the entire response to Discord
	discord.ChannelMessageSend(message.ChannelID, "Integrated Address Response:\n```\n"+outputMessage+"```")
}
