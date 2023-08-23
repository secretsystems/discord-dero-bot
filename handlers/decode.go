// handlers/integrated_address.go

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleIntegratedAddress(discord *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	// fmt.Printf("CONTENT: %s", content)
	if content == "!decode" {
		discord.ChannelMessageSend(message.ChannelID, "To decode an integrated address: `!decode <integrated_address>`")
		return

	} else if strings.HasPrefix(content, "!decode ") {

		userInput := strings.TrimPrefix(message.Content, "!decode ")
		// log.Printf("User Input: " + userInput)

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

		// Retrieve IP address, wallet port, and node port from environment variables
		ip := os.Getenv("DERO_SERVER_IP")
		derodPort := os.Getenv("DERO_WALLET_PORT")

		// Construct the URL using the retrieved IP address and wallet port
		url := fmt.Sprintf("http://%s:%s/json_rpc", ip, derodPort)

		// Define request for node
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		// Retrieve authentication credentials from environment variables
		username := os.Getenv("USER")
		password := os.Getenv("PASS")

		// Set basic authentication for the request
		request.SetBasicAuth(username, password)
		request.Header.Set("Content-type", "application/json")
		// fmt.Println("\nRequest: ", request)

		client := http.DefaultClient
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Error sending HTTP Post request: %v", err)
			return
		}
		defer response.Body.Close()

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

		// Print the entire mapResponse map
		var outputMessage string
		for key, value := range mapResponse {
			formattedValue, _ := json.MarshalIndent(value, "", "  ")
			outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
		}

		// Send the entire response to Discord
		discord.ChannelMessageSend(message.ChannelID, "Integrated Address Response:\n```\n"+outputMessage+"```")
	}
}
