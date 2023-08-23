package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

func HandleGetInfoDerod(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Define JSON struct
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "DERO.GetInfo",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		return
	}

	// Retrieve IP address, wallet port, and node port from environment variables
	ip := os.Getenv("DERO_SERVER_IP")
	derodPort := os.Getenv("DERO_NODE_PORT")

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
	request.Header.Set("Content-Type", "application/json")

	// fmt.Println("\nRequest: ", request)
	client := http.DefaultClient

	// Read response body
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Post request: %v", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code indicates authorization required
	if response.StatusCode == http.StatusUnauthorized {
		log.Printf("Authorization required for the request")
		return
	}

	responseBody, _ := io.ReadAll(response.Body)
	// log.Printf("Response Body: %v", string(responseBody))

	// Check if the response body is not empty
	if len(responseBody) > 0 {
		var mapResponse map[string]interface{}
		err = json.Unmarshal(responseBody, &mapResponse)
		if err != nil {
			log.Printf("Error decoding response JSON: %v", err)
			return
		}

		// Format response for printing
		var outputMessage string
		for key, value := range mapResponse {
			formattedValue, _ := json.MarshalIndent(value, "", " ")
			outputMessage += fmt.Sprintf("%s: %s\n", key, formattedValue)
		}

		// Send the entire response to Discord
		discord.ChannelMessageSend(message.ChannelID, "Node Info:\n```\n"+outputMessage+"```")
	} else {
		log.Printf("Empty response body")
	}
}
