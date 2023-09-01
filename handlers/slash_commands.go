// commands/slash_commands.go
package handlers

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commandIDs = make(map[string]string, len(Commands))
)

// Register slash commands
func RegisterSlashCommands(session *discordgo.Session, appID, guildID string) {
	log.Println("Adding Registering Slash Commands")
	for _, command := range Commands {
		registeredCommands, err := session.ApplicationCommandCreate(appID, guildID, &command)
		if err != nil {
			log.Fatalf("Cannot create %v slash command: %v", command.Name, err)
		}
		log.Printf("Registered Slash Commands: %v", command.Name)
		commandIDs[registeredCommands.ID] = registeredCommands.Name
		time.Sleep(2)
	}
}

func Cleanup(session *discordgo.Session, appID, guildID string) {
	log.Println("Adding Cleaning up")

	for id, name := range commandIDs {
		err := session.ApplicationCommandDelete(appID, guildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
		time.Sleep(2)
	}
}
