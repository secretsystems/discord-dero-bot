package handlers

import (
	"discord-dero-bot/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []discordgo.ApplicationCommand{
		{
			Name:        "encode",
			Description: "Encode Integrated Address",
		},
		{
			Name:        "trade-dero-xmr",
			Description: "Trade DERO-XMR",
		},
		{
			Name:        "decode",
			Description: "Decode Integrated Address",
		},
		{
			Name:        "giftbox",
			Description: "Get a DERO gift box!",
		},
	}
)
var (
	componentsHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string){
		"fd_yes": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Use the following addr in your DERO wallet to obtain instructions on how to sell your DERO for XMR ```deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7zat9hh2grpwfjjqcn4095kueeqg3z4yneqwa5hg6pqtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezq7rdwgs8gunpv3jjqctyv3ezqar0yp3x2gryv4kxjan9wfjkggr5dus8jmm4wgs8wctvd3jhgcjy25vs2wtzfe2sqcjk25pqrm2pm2```",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"fd_no": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Use the following addr in your DERO wallet to obtain instructions on how to buy DERO with XMR ```deroi1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqdyvfp4x7z7t9hh2grpwfjjqarjv9jxjmn8ypzy25j0ypnx7u3qtpx4ygp68gsyxmmdwpkx2ar9yp68sgrxdaezqarjv9jx2grfdeehgun4vd68xgr5dusxyefqv3jkc6tkv4ex2epqw3hjq7t0w4ezqampd3kx2arzg323j89rvf892qrz2e2sygeaw9a```",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}

	commandsHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string){
		"trade-dero-xmr": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "DERO-XMR is trading at: " + utils.ExchangeRateString() + "\nWould you like to trade? \nTrades have a fee of 1%",
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
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
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "DERO-XMR chart",
									Style:    discordgo.LinkButton,
									Disabled: false,
									URL:      "https://www.tradingview.com/chart/XAuuVNP7/",
								},
								discordgo.Button{
									Label:    "Walkthru",
									Style:    discordgo.LinkButton,
									Disabled: false,
									URL:      "https://youtu.be/x_EZ3BdpyyY",
								},
								discordgo.Button{
									Label:    "Github Repo",
									Style:    discordgo.LinkButton,
									Disabled: false,
									URL:      "https://github.com/secretnamebasis/dero-xmr-swap",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"encode": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "encode_" + interaction.Interaction.Member.User.ID,
					Title:    "Encode a DERO Integrated Address",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "address",
								Label:       "Address",
								Style:       discordgo.TextInputShort,
								Placeholder: "dero1q wallet address",
								Required:    true,
								MaxLength:   66,
								MinLength:   66,
							},
						},
						},
						discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "amount",
								Label:       "Amount; minimum 2 DERI",
								Style:       discordgo.TextInputShort,
								Placeholder: "How much in atomic units? 1 DERO = 100000",
								Required:    true,
								MaxLength:   64,
								MinLength:   1,
							},
						},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "comment",
									Label:     "What would you include in your encoding?",
									Style:     discordgo.TextInputParagraph,
									Required:  false,
									MaxLength: 128,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "destionation",
									Label:     "What port you want to send this too?, ex 1337",
									Style:     discordgo.TextInputShort,
									Required:  false,
									MaxLength: 128,
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"giftbox": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "giftbox_" + interaction.Interaction.Member.User.ID,
					Title:    "Purchase a DERO Gift Box",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "color",
								Label:       "Shirt Color",
								Style:       discordgo.TextInputShort,
								Placeholder: "black or white?",
								Required:    true,
							},
						},
						},
						discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "size",
								Label:       "Shirt Size",
								Style:       discordgo.TextInputShort,
								Placeholder: "What size fits you: S,M,L,XL,XXL,XXXL",
								Required:    true,
							},
						},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "address",
									Label:     "Shipping Address?",
									Style:     discordgo.TextInputParagraph,
									Required:  false,
									MaxLength: 80,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "Contact Info",
									Label:     "If we need to reach you?",
									Style:     discordgo.TextInputShort,
									Required:  false,
									MaxLength: 128,
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"decode": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, AppID, GuildID string) {
			err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "decode_" + interaction.Interaction.Member.User.ID,
					Title:    "Decode DERO Integrated Addr",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "integrated_address",
								Label:       "Integrated Address",
								Style:       discordgo.TextInputShort,
								Placeholder: "integrated wallet address",
								Required:    true,
								// MaxLength:   500,
								// MinLength:   66,
							},
						},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)
