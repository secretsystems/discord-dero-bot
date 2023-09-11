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

var chatGPTAPI string

const (
	secretMembersRoleID = "1057328486211145810"
	openAIURL           = "https://api.openai.com/v1/chat/completions"
)

func init() {
	chatGPTAPI = os.Getenv("OPEN_AI_TOKEN")
	if chatGPTAPI == "" {
		log.Println("OpenAI API token not found in environment")
	}
}

func hasSecretMembersRole(member *discordgo.Member) bool {
	if member == nil {
		return false
	}

	for _, roleID := range member.Roles {
		if roleID == secretMembersRoleID {
			return true
		}
	}
	return false
}

type ChatPayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func preparePayload(userInput string) ([]byte, error) {
	payload := ChatPayload{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: userInput},
		},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	return json.Marshal(payload)
}

func makeOpenAIRequest(payload []byte) ([]byte, error) {
	if chatGPTAPI == "" {
		return nil, fmt.Errorf("OpenAI API token not found")
	}

	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+chatGPTAPI)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func HandleChat(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if the message is sent in a DM channel
	if message.GuildID == "" {
		// This is a DM channel
		session.ChannelMessageSend(message.ChannelID, "You can't use the `!bot` command in DMs.")
		return
	}

	if !hasSecretMembersRole(message.Member) {
		session.ChannelMessageSend(message.ChannelID, "You don't have permission to use this command.\nTo gain permission, please consider becoming a `@secret-member`")
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!bot ")
	session.ChannelMessageSend(message.ChannelID, "Bot is processing your request:")

	userInput = userInput + " . Keep your response less than 1337 characters. Your max_tokens limit is 200"

	payload, err := preparePayload(userInput)
	if err != nil {
		log.Printf("Error encoding payload: %v", err)
		return
	}

	respBody, err := makeOpenAIRequest(payload)
	if err != nil {
		log.Printf("Error making OpenAI request: %v", err)
		return
	}

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

	err = json.Unmarshal(respBody, &chatResponse)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return
	}

	if len(chatResponse.Choices) > 0 {
		responseContent := chatResponse.Choices[0].Message.Content
		session.ChannelMessageSend(message.ChannelID, responseContent)
	}
}
