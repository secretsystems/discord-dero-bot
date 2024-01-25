package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/exports"
)

var (
	count = 10
)

func HandleMarketsRequest(session *discordgo.Session, message *discordgo.MessageCreate) {

	log.Printf("Initiating GET request to %s", exports.TradeOgreMarketsURL)

	// Create a GET request
	request, err := http.NewRequest("GET", exports.TradeOgreMarketsURL, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending HTTP GET request: %v", err)
		return
	}
	defer response.Body.Close()

	log.Println("Received response from TradeOgre API")

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	var marketData []map[string]map[string]string
	err = json.Unmarshal(responseBody, &marketData)
	if err != nil {
		log.Printf("Error decoding response JSON: %v", err)
		log.Printf("Response JSON: %s", string(responseBody))
		return
	}

	log.Println("Successfully decoded response JSON")

	// Extract and sort market pairs based on volume
	usdtPairs, btcPairs := extractAndSortPairs(marketData)

	// Create formatted lists for top 5 USDT pairs and top 5 BTC pairs

	log.Println("Formatted top pairs by volume")
	countStr := strconv.Itoa(count)
	// Combine the formatted pairs into one message
	combinedMessage := "Top " + countStr + " BTC Pairs TradeOgre.com:```" +
		formatPairs(btcPairs, count) + "```" +
		"Top " + countStr + " USDT Pairs TradeOgre.com:```" +
		formatPairs(usdtPairs, count) + "```"

	// Send the combined message to Discord
	session.ChannelMessageSend(message.ChannelID, combinedMessage)
}

func extractAndSortPairs(marketData []map[string]map[string]string) (usdtPairs, btcPairs []map[string]map[string]string) {
	// Separate pairs into USDT and BTC pairs
	for _, pairData := range marketData {
		pairName := getKey(pairData)
		if strings.HasSuffix(pairName, "-USDT") {
			usdtPairs = append(usdtPairs, pairData)
		} else if strings.HasSuffix(pairName, "-BTC") {
			btcPairs = append(btcPairs, pairData)
		}
	}

	// Sort pairs based on volume (converted to float64)
	sort.Slice(usdtPairs, func(i, j int) bool {
		return compareVolumes(usdtPairs[i], usdtPairs[j])
	})
	sort.Slice(btcPairs, func(i, j int) bool {
		return compareVolumes(btcPairs[i], btcPairs[j])
	})

	return usdtPairs, btcPairs
}

func compareVolumes(pair1, pair2 map[string]map[string]string) bool {
	volume1, err1 := strconv.ParseFloat(pair1[getKey(pair1)]["volume"], 64)
	volume2, err2 := strconv.ParseFloat(pair2[getKey(pair2)]["volume"], 64)

	if err1 != nil || err2 != nil {
		// Handle conversion errors (you may choose to log or ignore)
		log.Printf("Error converting volume to float64 for pair %s or %s", getKey(pair1), getKey(pair2))
		return false
	}

	return volume1 > volume2
}

func formatPairs(pairs []map[string]map[string]string, count int) string {
	var formattedPairs strings.Builder
	for i, pair := range pairs {
		if i >= count {
			break // Stop when the specified count is reached
		}
		formattedPairs.WriteString(formatPairDetails(pair))
	}
	return formattedPairs.String()
}

func formatPairDetails(pair map[string]map[string]string) string {
	pairName := getKey(pair)
	details := pair[pairName]
	return fmt.Sprintf("[%s]: Volume: %s, Price: %s\n",
		pairName, details["volume"], details["price"])
}

func getKey(pair map[string]map[string]string) string {
	for key := range pair {
		return key
	}
	return ""
}
