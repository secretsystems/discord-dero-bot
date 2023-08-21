// utils/dero/parse_transfers.go
package dero

import (
	"encoding/json"
	"fmt"
	"time"
)

// TransferEntry represents a transfer entry structure
type TransferEntry struct {
	Height          uint64      `json:"height"`
	TopoHeight      int64       `json:"topoheight"`
	BlockHash       string      `json:"blockhash"`
	MinerReward     uint64      `json:"minerreward"`
	TPoS            int         `json:"tpos"` // Position within the block; -1 for coinbase
	Pos             int         `json:"pos"`  // Position within the transaction
	Coinbase        bool        `json:"coinbase"`
	Incoming        bool        `json:"incoming"`
	TxID            string      `json:"txid"`
	Destination     string      `json:"destination"`
	Burn            uint64      `json:"burn,omitempty"`
	Amount          uint64      `json:"amount"`
	Fees            uint64      `json:"fees"`
	Proof           string      `json:"proof"` // Can be used to prove if available
	Status          byte        `json:"status"`
	Time            time.Time   `json:"time"`
	EWData          string      `json:"ewdata"` // Encrypted wallet balance at that point in time
	Data            []byte      `json:"data"`   // Data is the entire decrypted dump
	PayloadType     byte        `json:"payloadtype"`
	Payload         []byte      `json:"payload"`
	PayloadError    string      `json:"payloaderror,omitempty"`
	PayloadRPC      []Arguments `json:"payload_rpc,omitempty"`
	Sender          string      `json:"sender"`
	DestinationPort uint64      `json:"dstport"`
	SourcePort      uint64      `json:"srcport"`
}

// Arguments is a placeholder type for the Payload_RPC field
type Arguments struct {
	Name     string      `json:"name"`
	Datatype string      `json:"datatype"`
	Value    interface{} `json:"value"`
}

type ResponseData struct {
	Result struct {
		Entries []TransferEntry `json:"entries"`
	} `json:"result"`
}

func ParseTransfersResponse(responseBody []byte) ([]TransferEntry, error) {
	var responseData ResponseData

	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return responseData.Result.Entries, nil
}
