// handlers/echo_handler.go

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
	fmt.Println("\nResponse Body:", string(responseBody))

	var httpResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &httpResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return
	}

	// Print the entire httpResponse map
	fmt.Println("\nhttpResponse:", httpResponse)

	// Print individual fields if present
	if result, ok := httpResponse["result"].(map[string]interface{}); ok {
		fmt.Println("\nResult:", result)

		// If "payload_rpc" field is present
		if payloadRPC, ok := result["payload_rpc"].([]interface{}); ok {
			for _, payloadMap := range payloadRPC {
				if payload, ok := payloadMap.(map[string]interface{}); ok {
					if value, ok := payload["value"].(string); ok {
						if strings.HasPrefix(value, "You are trading") {
							output := strings.TrimPrefix(value, "You are trading")
							fmt.Println("Output:", output)
							log.Printf("Echo response: " + output)
							discord.ChannelMessageSend(message.ChannelID, "Integrated Address: "+output)
							return
						}
					}
				}
			}
		}
	}
}
