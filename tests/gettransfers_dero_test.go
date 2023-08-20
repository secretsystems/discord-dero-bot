// tests/gettransfers_dero_test.go
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fuck_you.com/utils"
	"github.com/stretchr/testify/assert"
)

func TestFetchDeroTransfers(t *testing.T) {
	// Create a mock DeroServer
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockData := map[string]interface{}{
			"result": "mock data",
		}
		mockDataJSON, _ := json.Marshal(mockData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockDataJSON)
	}))
	defer mockServer.Close()

	// Create a new DeroServer instance
	server := utils.NewDeroServer(mockServer.URL, "user", "pass")

	// Call the FetchTransfers method
	transfersData, err := server.FetchTransfers()
	if err != nil {
		t.Fatalf("Error fetching Dero transfers: %v", err)
	}

	// Define the expected JSON result
	want := `{"result": "mock data"}`

	// Compare the expected result with the actual result
	assert.JSONEq(t, want, string(transfersData), "Expected and actual JSON data do not match.")
}
