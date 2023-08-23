// commands/slash_commands.go
package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Register slash commands
func RegisterSlashCommands(discordSession *discordgo.Session, AppID, GuildID string) {
	_, err := discordSession.ApplicationCommandCreate(AppID, GuildID, &discordgo.ApplicationCommand{
		Name:        "trade-dero-xmr",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
}
