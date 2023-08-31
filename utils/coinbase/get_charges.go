package coinbase

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetCharges() {
	coinbaseAPIToken := os.Getenv("COINBASE_API_TOKEN")
	url := "https://api.commerce.coinbase.com/checkouts/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-CC-Api-Key", coinbaseAPIToken) // Set the Coinbase API token
	req.Header.Add("X-CC-Version", "2018-03-22")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
