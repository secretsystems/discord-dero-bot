package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMarketsRequest(discord *discordgo.Session, message *discordgo.MessageCreate) {
	url := "https://tradeogre.com/api/v1/markets"

	// Create a GET request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP Get request: %v", err)
		return
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	log.Printf("Response Body: %v", string(responseBody))

	var marketData []map[string]map[string]string
	err = json.Unmarshal(responseBody, &marketData)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		return
	}

	pairs := []string{}
	for _, pairData := range marketData {
		for pairName := range pairData {
			pairs = append(pairs, pairName)
		}
	}

	// Custom sorting function
	sort.SliceStable(pairs, func(i, j int) bool {
		quoteI := strings.SplitN(pairs[i], "-", 2)[1]
		quoteJ := strings.SplitN(pairs[j], "-", 2)[1]
		return quoteI < quoteJ
	})

	// Create a map to store pairs grouped by quotes
	pairsByQuote := make(map[string][]string)
	for _, pair := range pairs {
		quote := strings.SplitN(pair, "-", 2)[1]
		pairsByQuote[quote] = append(pairsByQuote[quote], pair)
	}

	// Create a formatted list of sorted pairs with headers
	var formattedPairs strings.Builder
	for quote, quotePairs := range pairsByQuote {
		formattedPairs.WriteString(fmt.Sprintf("%s:\n%s\n\n", quote, strings.Join(quotePairs, " ")))
	}

	// Send the sorted pairs to Discord
	discord.ChannelMessageSend(message.ChannelID, "Sorted Market Pairs:\n```\n"+formattedPairs.String()+"```")
}
