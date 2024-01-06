package dero

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/SixofClubsss/Grokked/grok"
	drpc "github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/rpc"
)

const SIX_OF_CLUBS_DISCORD_USER_ID = "767476835558096928"

var userAddressMap = make(map[string]string)
var userMappingsMutex sync.Mutex

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
			mappedAddress := getUserAddress(lbl)

			if mappedAddress != "" {
				lbl = "<@" + mappedAddress + ">"
			}
			return fmt.Sprintf("# Grok is...\n> %s \n> Time: ```%s```\n> **dApp by: https://dreamdapps.io**", lbl, left)
		}
	}

	return "I am not sure? contact the dev <@" + SIX_OF_CLUBS_DISCORD_USER_ID + ">"
}

func getUserAddress(lbl string) string {
	userMappingsMutex.Lock()
	defer userMappingsMutex.Unlock()
	return userAddressMap[lbl]
}
