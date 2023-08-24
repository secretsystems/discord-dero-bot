// commands/slash_commands.go
package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commandIDs = make(map[string]string, len(Commands))
)

// Register slash commands
func RegisterSlashCommands(discord *discordgo.Session, AppID, GuildID string) {
	for _, command := range Commands {
		registeredCommands, err := discord.ApplicationCommandCreate(AppID, GuildID, &command)
		if err != nil {
			log.Fatalf("Cannot create %v slash command: %v", &command, err)
		}
		commandIDs[registeredCommands.ID] = registeredCommands.Name
	}
}

func Cleanup(discord *discordgo.Session, AppID, GuildID string) {
	for id, name := range commandIDs {
		err := discord.ApplicationCommandDelete(AppID, GuildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}
