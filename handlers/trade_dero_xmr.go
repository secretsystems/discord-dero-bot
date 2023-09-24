package handlers

import (
	"discord-dero-bot/utils"

	"github.com/bwmarrin/discordgo"
)

func handleFdYes(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "Use the following address in your DERO wallet to obtain instructions on how to Buy DERO with XMR```deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7zat9hh2grpwfjjqcn4095kueeqg3z4yneqwa5hg6pqtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezq7rdwgs8gunpv3jjqctyv3ezqar0yp3x2gryv4kxjan9wfjkggr5dus8jmm4wgs8wctvd3jhgcjy25vs2wtzfe2sqcjk25pqrm2pm2```"
	RespondWithMessage(session, interaction, message)
}
func handleFdNo(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "Use the following address in your DERO wallet to obtain instructions on how to sell DERO for XMR```deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7z7t9hh2grpwfjjqarjv9jxjmn8ypzy25j0ypnx7u3qtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezqarjv9jx2grfdeehgun4vd68xgr5dusxyefqv3jkc6tkv4ex2epqw3hjq7t0w4ezqampd3kx2arzg323j89rvf892qrz2e2sygeaw9a```"
	RespondWithMessage(session, interaction, message)
}

func handleTradeDeroXmrComponent(session *discordgo.Session, interaction *discordgo.InteractionCreate, appID string) {
	message := "DERO-XMR is trading at: " + utils.DeroXmrExchangeRateString() + "\nWould you like to trade? \nTrades have a fee of 1%"

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "BUY DERO with XMR",
				Style:    discordgo.SuccessButton,
				Disabled: false,
				CustomID: "fd_yes",
			},
			discordgo.Button{
				Label:    "SELL DERO for XMR",
				Style:    discordgo.DangerButton,
				Disabled: false,
				CustomID: "fd_no",
			},
		}},
	}

	urlButtons := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "DERO-XMR chart",
				Style:    discordgo.LinkButton,
				Disabled: false,
				URL:      "https://www.tradingview.com/chart/XAuuVNP7/",
			},
			discordgo.Button{
				Label:    "Github Repo",
				Style:    discordgo.LinkButton,
				Disabled: false,
				URL:      "https://github.com/secretnamebasis/dero-xmr-swap",
			},
		}},
	}

	videoButtons := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Walkthru: Buy DERO",
				Style:    discordgo.LinkButton,
				Disabled: false,
				URL:      "https://youtu.be/OGuV7jSAccE",
			},
			discordgo.Button{
				Label:    "Walkthru: Sell DERO",
				Style:    discordgo.LinkButton,
				Disabled: false,
				URL:      "https://youtu.be/RLeN03QC6jE",
			},
		}},
	}

	// Create the components slice by combining the button slices
	components := append(buttons, urlButtons...)
	components = append(components, videoButtons...)

	respondWithMessageAndComponents(session, interaction, message, components)
}
