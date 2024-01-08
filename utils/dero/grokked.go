package dero

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SixofClubsss/Grokked/grok"
	drpc "github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/rpc"
)

const SIX_OF_CLUBS_DISCORD_USER_ID = "767476835558096928"

var (
	userMappings      = make(map[string]string)
	addressMappings   = make(map[string]string)
	userMappingsMutex sync.Mutex
)

// Get single string key result from scid
func GetStringKey(scid, key string) interface{} {
	client, ctx, cancel := drpc.SetDaemonClient(fmt.Sprintf("%s:%s", DeroServerIP, DeroServerPort))
	defer cancel()

	var result *rpc.GetSC_Result
	params := rpc.GetSC_Params{
		SCID:      scid,
		Code:      false,
		Variables: true,
	}

	if err := client.CallFor(ctx, &result, "DERO.GetSC", params); err != nil {
		log.Println("[GetStringKey]", err)
		return nil
	}

	return result.VariableStringKeys[key]
}

// Get single uint64 key result from scid
func GetUintKey(scid, key string) interface{} {
	client, ctx, cancel := drpc.SetDaemonClient(fmt.Sprintf("%s:%s", DeroServerIP, DeroServerPort))
	defer cancel()

	var result *rpc.GetSC_Result
	params := rpc.GetSC_Params{
		SCID:      scid,
		Code:      false,
		Variables: true,
	}

	if err := client.CallFor(ctx, &result, "DERO.GetSC", params); err != nil {
		log.Println("[GetUintKey]", err)
		return nil
	}

	return result.VariableUint64Keys[drpc.StringToUint64(key)]
}

// Get the current Grok and time frame on GROKSCID
func GetGrok() string {
	loadUserMap()
	if f, ok := GetStringKey(grok.GROKSCID, "grok").(float64); ok {
		if raw, ok := GetUintKey(grok.GROKSCID, strconv.FormatFloat(f, 'f', 0, 64)).(string); ok {
			left := "I am not sure how much time is left? contact the dev <@" + SIX_OF_CLUBS_DISCORD_USER_ID + ">"
			now := float64(time.Now().Unix())
			if last, ok := GetStringKey(grok.GROKSCID, "last").(float64); ok {
				if dur, ok := GetStringKey(grok.GROKSCID, "duration").(float64); ok {
					tf := last + dur
					if now < tf {
						left = fmt.Sprintf("%.0f minutes left to pass", (tf-now)/60)
					} else if tf != 0 {
						left = fmt.Sprintf("%.0f minutes past", (now-tf)/60)
					}
				}
			}

			lbl := drpc.DeroAddressFromKey(raw)
			mappedID := getUserID(lbl)

			if mappedID != "" {
				lbl = "<@" + mappedID + ">"
			}
			return fmt.Sprintf("# Grok is...\n> %s \n> Time: ```%s```\n> **dApp by: https://dreamdapps.io**", lbl, left)
		}
	}

	return "I am not sure? contact the dev <@" + SIX_OF_CLUBS_DISCORD_USER_ID + ">"
}

func getUserID(lbl string) string {
	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()
	return addressMappings[lbl]
}

func loadUserMap() error {
	data, err := os.ReadFile("userMappings.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &userMappings)
	if err != nil {
		return err
	}

	// Build the reverse map for address-based lookup
	for userID, address := range userMappings {
		addressMappings[address] = userID
	}

	// Ensure keys in the maps are trimmed and in lowercase
	// Modify the content of userMappings
	for k, v := range userMappings {
		delete(userMappings, k)
		userMappings[strings.TrimSpace(strings.ToLower(k))] = strings.TrimSpace(strings.ToLower(v))
	}

	return nil
}
