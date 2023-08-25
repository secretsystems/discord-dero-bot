package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleChat(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if the user has the required role
	hasSecretMembersRole := false
	for _, roleID := range message.Member.Roles {
		if roleID == "1057328486211145810" { // Change this to the actual role ID
			hasSecretMembersRole = true
			break
		}
	}

	if !hasSecretMembersRole {
		// The user doesn't have the required role, return or send an error message
		discord.ChannelMessageSend(message.ChannelID, "You don't have permission to use this command.")
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!bot ")
	discord.ChannelMessageSend(message.ChannelID, "Bot is processing your request:")

	// fmt.Printf(userInput)
	// Prepare the request payload
	payload := struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"max_tokens"`
	}{
		Model: "gpt-3.5-turbo",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: userInput},
		},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding payload: %v", err)
		return
	}

	// Retrieve the OpenAI API token from the environment
	apiToken := os.Getenv("OPEN_AI_TOKEN")
	// fmt.Printf(apiToken)
	if apiToken == "" {
		log.Println("OpenAI API token not found in environment")
		return
	}

	// Make the API request to OpenAI
	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}
	// fmt.Printf("requst: %v\n", req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// fmt.Printf("Response Body: %s\n", respBody)

	var chatResponse struct {
		ID      string `json:"id"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
			Index        int    `json:"index"`
		} `json:"choices"`
	}

	err = json.Unmarshal(respBody, &chatResponse) // Use Unmarshal instead of NewDecoder
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return
	}

	// Send the response to Discord
	if len(chatResponse.Choices) > 0 {
		responseContent := chatResponse.Choices[0].Message.Content
		discord.ChannelMessageSend(message.ChannelID, responseContent)
	}
}
