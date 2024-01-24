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
	buyDero              = "deroi1qykz2fqtptcnvmr65042jwljpwmglnezax4wms5w4htat20vzsdauq9yvfp4x7zd23exzer9ypzy25j0ypnx7u3qtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezq7rdwgsxzerywgs8gmeqvfjjqer9d35hvetjv4jzqar0ypuk7atjypmkzmrvv46xy3z4ryznjcjw25qxy4j4qg94lam6"
	sellDero             = "deroi1qykz2fqtptcnvmr65042jwljpwmglnezax4wms5w4htat20vzsdauq9yvfp4x7zv23exzerfdenjqkzd2gsxvmmjypzy25j0yqar5gzrdakhqmr9w3jjqarcypnx7u3qw3exzer9yp6x7gryv4kxjan9wgsxjmnxdus8gmeq09hh2u3qwaskcmr9w33yg4gerj3kynj4qp39v4gzwgckeh"
	emojiForBuy          = discordgo.ComponentEmoji{Name: "🟢"}
	emojiForSell         = discordgo.ComponentEmoji{Name: "🔴"}
	emojiForYouTube      = discordgo.ComponentEmoji{Name: "🎥"}
	emojiForGitHub       = discordgo.ComponentEmoji{Name: "💻"}
	emojiForTradingView  = discordgo.ComponentEmoji{Name: "📈"}
	emojiForConfirmation = discordgo.ComponentEmoji{Name: "✅"}
)
var confirm = false

func handleFdWalkthrus(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "# Walktrhus\n" +
		"## Android:\n" +
		"Buy DERO on with ENGRAM: https://youtu.be/PsGW29X_ze8\n" +
		"## macOS:\n" +
		"Buy DERO on macOS with CLI: https://youtu.be/OGuV7jSAccE\n" +
		"Sell DERO on macOS with CLI: https://youtu.be/RLeN03QC6jE\n" +
		""

	RespondWithMessage(session, interaction, message)
}

func handleFdConfirmBuy(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "An integrated address has been populated, please click confirm to continue\n" +
		"> Note: After you click confirm, you will see an integrated address. Use this address " +
		"in your DERO wallet of choice and complete the micro transfer. Once completed you will" +
		" receive instructions in your DERO wallet transaction history for how **buy DERO for XMR**"
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Confirm",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					CustomID: "fd_yes",
					Emoji:    emojiForConfirmation,
				},
				discordgo.Button{
					Label:    "Walkthrus",
					Style:    discordgo.SecondaryButton,
					Disabled: false,
					CustomID: "fd_walkthru",
					Emoji:    emojiForYouTube,
				},
			},
		},
	}
	respondWithMessageAndComponents(session, interaction, message, buttons)
}
func handleFdConfirmSell(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "An integrated address has been populated, please click confirm to continue\n" +
		"> Note: After you click confirm, you will see an integrated address. Use this address " +
		"in your DERO wallet of choice and complete the micro transfer. Once completed you will" +
		" receive instructions in your DERO wallet transaction history for how **sell DERO for XMR**"
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Confirm",
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					CustomID: "fd_no",
					Emoji:    emojiForConfirmation,
				},
				discordgo.Button{
					Label:    "Walkthrus",
					Style:    discordgo.SecondaryButton,
					Disabled: false,
					CustomID: "fd_walkthru",
					Emoji:    emojiForYouTube,
				},
			},
		},
	}
	respondWithMessageAndComponents(session, interaction, message, buttons)
}

func handleFdYes(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// message := "Use the following address in your DERO wallet to obtain instructions on how to Buy DERO with XMR```" + buyDero + "```"

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Wallet: CLI",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/deroproject/derohe/releases/latest",
					Emoji:    emojiForGitHub,
				},
				discordgo.Button{
					Label:    "Wallet: ENGRAM",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/DEROFDN/Engram/releases/latest",
					Emoji:    emojiForGitHub,
				},
			},
		},
	}

	respondWithMessageAndComponents(session, interaction, buyDero, buttons)
}
func handleFdNo(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	// message := "Use the following address in your DERO wallet to obtain instructions on how to sell DERO for XMR```" + sellDero + "```"
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Wallet: CLI",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/deroproject/derohe/releases/latest",
					Emoji:    emojiForGitHub,
				},
				discordgo.Button{
					Label:    "Wallet: ENGRAM",
					Style:    discordgo.LinkButton,
					Disabled: false,
					URL:      "https://github.com/DEROFDN/Engram/releases/latest",
					Emoji:    emojiForGitHub,
				},
			},
		},
	}

	respondWithMessageAndComponents(session, interaction, sellDero, buttons)
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
					CustomID: "fd_confirm_buy",
					Emoji:    emojiForBuy,
				},
				discordgo.Button{
					Label:    "SELL DERO - XMR",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: "fd_confirm_sell",
					Emoji:    emojiForSell,
				},
			},
		},
	}

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
