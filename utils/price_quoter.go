package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type Response struct {
	Success      bool
	Initialprice string
	Price        string
	High         string
	Low          string
	Volume       string
	Bid          string
	Ask          string
}

func getJson(pair string, target interface{}) error {

	resp, err := client.Get("https://tradeogre.com/api/v1/ticker/" + pair)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// well I want the price from trade ogre for dero-usdt
func getAsk(pair string) float64 {
	result := Response{}
	err := getJson(pair, &result)
	if err != nil {
		log.Fatalf("Failed to fetch %s: %v", pair, err)
	}

	parsed, perr := strconv.ParseFloat(result.Ask, 64)
	if perr != nil {
		log.Fatalf("failed to parse float: %v", perr)
	}
	return parsed
}

func ExchangeRate() float64 {
	dero := getAsk("dero-usdt")
	fmt.Printf("The Price of DERO: %v", dero)
	return dero
}

func ExchangeRateInt(f float64) int {
	return int(f)
}

func ExchangeRateString() string {
	return fmt.Sprintf("%f", ExchangeRate())
}
