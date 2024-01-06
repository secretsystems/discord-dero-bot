package handlers

import (
	"fmt"
	"strconv"

	"github.com/secretsystems/discord-dero-bot/utils"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
	"github.com/secretsystems/discord-dero-bot/utils/monero"

	"github.com/bwmarrin/discordgo"
)

var (
	buyDero             = "deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7zd23exzer9ypzy25j0ypnx7u3qtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezq7rdwgsxzerywgs8gmeqvfjjqer9d35hvetjv4jzqar0ypuk7atjypmkzmrvv46xy3z4ryznjcjw25qxy4j4qgszhplu"
	sellDero            = "deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7zv23exzerfdenjqkzd2gsxvmmjypzy25j0yqar5gzrdakhqmr9w3jjqarcypnx7u3qw3exzer9yp6x7gryv4kxjan9wgsxjmnxdus8gmeq09hh2u3qwaskcmr9w33yg4gerj3kynj4qp39v4gzxtqv8j"
	emojiForBuy         = discordgo.ComponentEmoji{Name: "🟢"}
	emojiForSell        = discordgo.ComponentEmoji{Name: "🔴"}
	emojiForYouTube     = discordgo.ComponentEmoji{Name: "🎥"}
	emojiForGitHub      = discordgo.ComponentEmoji{Name: "💻"}
	emojiForTradingView = discordgo.ComponentEmoji{Name: "📈"}
)

func handleFdYes(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "Use the following address in your DERO wallet to obtain instructions on how to Buy DERO with XMR```" + buyDero + "```"

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Walkthru: Buy DERO",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://youtu.be/OGuV7jSAccE",
					Emoji:    emojiForYouTube,
				},
			},
		},
	}
	buttons = append(
		buttons,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Wallet: ENGRAM",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/DEROFDN/Engram/releases/latest",
					Emoji:    emojiForGitHub,
				},
			},
		},
	)
	respondWithMessageAndComponents(session, interaction, message, buttons)
}
func handleFdNo(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "Use the following address in your DERO wallet to obtain instructions on how to sell DERO for XMR```" + sellDero + "```"
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Walkthru: Sell DERO",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://youtu.be/RLeN03QC6jE",
					Emoji:    emojiForYouTube,
				},
			},
		},
	}
	buttons = append(
		buttons,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Wallet: ENGRAM",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/DEROFDN/Engram/releases/latest",
					Emoji:    emojiForGitHub,
				},
			},
		},
	)
	respondWithMessageAndComponents(session, interaction, message, buttons)
}

func handleTradeDeroXmrComponent(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	moneroBal, _ := monero.GetWalletBalance()
	deroBal, _ := dero.GetDeroWalletBalance()

	// Convert balances to the desired decimal precision
	deroBalStr := strconv.FormatFloat(float64(deroBal)/100000, 'f', 5, 64)
	moneroBalStr := strconv.FormatFloat(float64(moneroBal)/1000000000000, 'f', 12, 64)

	message := fmt.Sprintf(
		"DERO-XMR is trading at: %s\nThe `secret-wallet` has:\n> DERO %s\n> XMR %s\nWould you like to trade?\nTrades have a fee of 1%%",
		utils.DeroXmrExchangeRateString(), deroBalStr, moneroBalStr,
	)

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "BUY DERO - XMR",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: "fd_yes",
					Emoji:    emojiForBuy,
				},
			},
		},
	}
	buttons = append(
		buttons,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "SELL DERO - XMR",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: "fd_no",
					Emoji:    emojiForSell,
				},
			},
		},
	)
	buttons = append(
		buttons,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Github Repo",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/secretnamebasis/dero-xmr-swap",
					Emoji:    emojiForGitHub,
				},
			},
		},
	)
	buttons = append(
		buttons,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "DERO - XMR chart",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://www.tradingview.com/chart/XAuuVNP7/",
					Emoji:    emojiForTradingView,
				},
			},
		},
	)

	respondWithMessageAndComponents(session, interaction, message, buttons) // can be no more than five action rows high
}
