package dero

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Import the godotenv package
	"github.com/ybbus/jsonrpc"
)

var (
	// dero
	DeroServerIP     string
	DeroServerPort   string
	deroTipsPort     string
	deroWalletPort   string
	deroUser         string
	deroPass         string
	pongAmount       = "1337331"
	pongDir          string
	pongDB           string
	iAddressTextFile string
	DeroHttpClient   *http.Client
	DeroRpcClient    jsonrpc.RPCClient
)

type TransportWithBasicAuth struct {
	Username string
	Password string
	Base     http.RoundTripper
}

// RoundTrip implements the RoundTripper interface
func (t *TransportWithBasicAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Base.RoundTrip(req)
}

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//dero
	DeroServerIP = os.Getenv("DERO_SERVER_IP")
	deroWalletPort = os.Getenv("DERO_WALLET_PORT")
	deroTipsPort = os.Getenv("DERO_TIPS_PORT")
	DeroServerPort = os.Getenv("DERO_NODE_PORT")
	deroUser = os.Getenv("DERO_WALLET_USER")
	deroPass = os.Getenv("DERO_WALLET_PASS")
	homeDir, _ := os.UserHomeDir()
	pongDir = homeDir + "/dero-utils"
	pongDB = pongDir + "/" + pongAmount + ".sales.db"
	iAddressTextFile = pongDir + "/" + pongAmount + ".iaddress.txt"

	DeroHttpClient = &http.Client{
		Transport: &TransportWithBasicAuth{
			Username: deroUser,
			Password: deroPass,
			Base:     http.DefaultTransport,
		},
	}

	DeroRpcClient = jsonrpc.NewClientWithOpts(
		fmt.Sprintf("%s:%s", DeroServerIP, DeroServerPort),
		&jsonrpc.RPCClientOpts{
			HTTPClient: DeroHttpClient,
		},
	)

	if _, err := os.Stat(pongDir); os.IsNotExist(err) {
		err := os.Mkdir(pongDir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}
