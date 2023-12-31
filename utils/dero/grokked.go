package dero

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SixofClubsss/Grokked/grok"
	drpc "github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/rpc"
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
	if f, ok := GetStringKey(grok.GROKSCID, "grok").(float64); ok {
		if raw, ok := GetUintKey(grok.GROKSCID, strconv.FormatFloat(f, 'f', 0, 64)).(string); ok {
			left := "I am not sure how much time is left? contact the dev"
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

			return fmt.Sprintf("Grok is currently:\n\n%s\n\n%s", drpc.DeroAddressFromKey(raw), left)
		}
	}

	return "I am not sure? contact the dev"
}
