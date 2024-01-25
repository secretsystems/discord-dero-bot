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
	"github.com/secretsystems/discord-dero-bot/exports"
)

var chatGPTAPI string

func init() {
	chatGPTAPI = os.Getenv("OPEN_AI_TOKEN")
	if chatGPTAPI == "" {
		log.Println("OpenAI API token not found in the environment")
	}
}

func hasSecretMembersRole(session *discordgo.Session, guildID, roleID, userID string) bool {
	member, err := session.GuildMember(exports.SecretGuildID, userID)
	if err != nil {
		log.Printf("Error getting guild member: %v", err)
		return false
	}

	for _, memberRoleID := range member.Roles {
		if memberRoleID == roleID {
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
		Model:       "gpt-3.5-turbo",
		Messages:    []Message{{Role: "user", Content: userInput}},
		Temperature: 0.7,
		MaxTokens:   200,
	}

	return json.Marshal(payload)
}

func makeOpenAIRequest(payload []byte) ([]byte, error) {
	if chatGPTAPI == "" {
		return nil, fmt.Errorf("OpenAI API token not found")
	}

	req, err := http.NewRequest("POST", exports.OpenAIURL, bytes.NewBuffer(payload))
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
	if message.GuildID == "" {
		session.ChannelMessageSend(message.ChannelID, "You can't use the `!bot` command in DMs.")
		return
	}

	if !hasSecretMembersRole(session, message.GuildID, exports.SecretMembersRoleID, message.Author.ID) {
		session.ChannelMessageSend(message.ChannelID, "You don't have permission to use this command.\nTo gain permission, please consider becoming a `@secret-member` in https://discord.gg/BKdX9qHkgu")
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!bot ")
	session.ChannelMessageSend(message.ChannelID, "Bot is processing your request:")

	userInput = userInput + ". Keep your response less than 1337 characters. Your max_tokens limit is 200"

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
