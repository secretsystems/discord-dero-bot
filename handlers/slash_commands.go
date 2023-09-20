package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commandIDs = make(map[string]string, len(Commands))
)

type TokenBucket struct {
	tokens       int
	refillAmount int
	refillRate   time.Duration
	lastRefill   time.Time
	mutex        sync.Mutex
}

func NewTokenBucket(tokens, refillAmount int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:       tokens,
		refillAmount: refillAmount,
		refillRate:   refillRate,
		lastRefill:   time.Now(),
	}
}

func (tb *TokenBucket) Wait(n int) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	elapsed := time.Since(tb.lastRefill)
	tokensToRefill := int(elapsed / tb.refillRate)
	if tokensToRefill > tb.refillAmount {
		tokensToRefill = tb.refillAmount
	}

	tb.tokens += tokensToRefill
	tb.lastRefill = tb.lastRefill.Add(time.Duration(tokensToRefill) * tb.refillRate)

	for tb.tokens < n {
		time.Sleep(tb.refillRate)
		tb.tokens++
	}

	tb.tokens -= n
}

func RegisterSlashCommands(session *discordgo.Session, appID, guildID string, registrationBucket *TokenBucket) {
	log.Println("Starting registration of Slash Commands")

	for _, command := range Commands {
		registrationBucket.Wait(1)
		registeredCommands, err := session.ApplicationCommandCreate(appID, guildID, &command)
		if err != nil {
			log.Fatalf("Cannot create %v slash command: %v", command.Name, err)
		}
		log.Printf("Registered Slash Command: %v", command.Name)
		commandIDs[registeredCommands.ID] = registeredCommands.Name
	}

	log.Println("Finished registration of Slash Commands")
}

func Cleanup(session *discordgo.Session, appID, guildID string, cleanupBucket *TokenBucket) {
	log.Println("Starting cleanup of Slash Commands")

	// Fetch the current command IDs associated with the application and guild
	currentCommands, err := GetGuildSlashCommands(session, appID, guildID)
	if err != nil {
		log.Fatalf("Cannot fetch application commands for appID %v and guildID %v: %v", appID, guildID, err)
		return
	}

	// Create a set to store the command IDs for quick lookup
	commandIDSet := make(map[string]bool)
	for _, cmd := range currentCommands {
		commandIDSet[cmd.ID] = true
	}

	// Iterate over the commandIDs and delete the commands from Discord
	for id, name := range commandIDs {
		if _, exists := commandIDSet[id]; !exists {
			log.Printf("Command %v not found on Discord. Skipping cleanup.", name)
			continue
		}

		cleanupBucket.Wait(1)

		log.Printf("Attempting to clean up Slash Command: ID=%v, Name=%v", id, name)

		err := session.ApplicationCommandDelete(appID, guildID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		} else {
			log.Printf("Successfully deleted Slash Command: %v", name)
		}

		// Remove the command from the commandIDs map after successful deletion
		delete(commandIDs, id)
	}

	// Print the remaining command IDs
	log.Println("Remaining Command IDs:")
	for id := range commandIDSet {
		log.Println("Command ID:", id)

	}

	log.Println("Finished cleanup of Slash Commands")
}

func GetGuildSlashCommands(session *discordgo.Session, appID, guildID string) ([]*discordgo.ApplicationCommand, error) {
	commands, err := session.ApplicationCommands(appID, guildID)
	if err != nil {
		return nil, err
	}
	return commands, nil
}
