package main

import (

	"discord-dero-bot/bot"      
	"discord-dero-bot/handlers" 
	"discord-dero-bot/utils/coinbase"

	"discord-dero-bot/utils/dero"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	// discord
	botToken          = "MTE0MTEwMTU4ODcxMTIyNzQzNA.GUiXSz.URz9X36qAU152QdZA4jdNVk0tKixMSsaa4vSVc"
	guildID           = "1026670486358335579"
	appID             = "1141101588711227434"
	resultsChannel    = "1026670486895210618"
	registerationRole = "1147202657178624111"
	membershipRole    = "1147256770847314161"
	membershipAmount  = int((5 / float64(utils.ExchangeRate())) * 100000)
	// DERO
	serverIP      = "192.168.0.199"
	walletAddress = "dero1qyw4fl3dupcg5qlrcsvcedze507q9u67lxfpu8kgnzp04aq73yheqqg2ctjn4"
	walletPort    = "10103"
	serverPort    = "10102"
	user          = "user"
	pass          = "pass"
	homedir, _    = os.UserHomeDir()
	utilsdir      = homedir + "/dero-utils"
	amount        = "1337"
)

func init() {
	dero.InitializeDERO(
		serverIP,
		walletPort,
		serverPort,
		user,
		pass,
		homedir,
		utilsdir,
		membershipAmount,
	)

	if _, err := os.Stat(utilsdir); os.IsNotExist(err) {
		err := os.Mkdir(utilsdir, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}
}

func main() {
	// Initialize the bot
	bot, err := bot.NewBot(botToken) // Replace with the actual initialization function
	if err != nil {
		log.Fatalf("Error initializing Discord bot: %v", err)
	}

	bot.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening Discord bot connection: %v", err)
	}
	defer bot.Close()

	// Get the Discord session from the bot instance
	session := bot.GetDiscordSession()

	// Register interaction handlers

	handlers.AddHandlers(session, appID, guildID)
	handlers.AddModals(session, appID, guildID, resultsChannel)
	handlers.RegisterSlashCommands(session, appID, guildID)
	handlers.SetRegisteredRole(registerationRole)
	handlers.SetDesiredRole(membershipRole)
	handlers.SetDeroAddress(walletAddress)
	handlers.SetMembershipAmount(membershipAmount)
	handlers.SetResultsChannel(resultsChannel)

	// Call FetchAndParseTransfers function from the utils package
	log.Printf("Initializing DERO\n")
	transferEntries, err := dero.FetchAndParseTransfers()
	if err != nil {
		log.Printf("Error fetching and parsing transfers: %v", err)
	} else {
		// Process the fetched and parsed transfer entries
		log.Printf("Fetched and parsed %d transfer entries.\n", len(transferEntries))
	}

	log.Println("Bot is running. Press Ctrl+C to stop.")

	// Set up a channel to capture the Ctrl+C signal
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Wait for an interrupt signal to close the program
	<-channel
	handlers.Cleanup(session, appID, guildID)
	log.Println("Bot is cleaning up.")

}
