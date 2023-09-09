// commands/slash_commands.go
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

func RegisterSlashCommands(session *discordgo.Session, appID, guildID string, registrationBucket *TokenBucket) {
	log.Println("Starting registration of Slash Commands")

	for _, command := range Commands {
		registrationBucket.Wait(1) // Wait for a token to be available
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

	for id, name := range commandIDs {
		cleanupBucket.Wait(1) // Wait for a token to be available
		err := session.ApplicationCommandDelete(appID, guildID, id)
		log.Printf("Cleaning up Slash Command: %v\n", name)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

	log.Println("Finished cleanup of Slash Commands")
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

	// Calculate how many tokens need to be refilled
	elapsed := time.Since(tb.lastRefill)
	tokensToRefill := int(elapsed / tb.refillRate)
	if tokensToRefill > tb.refillAmount {
		tokensToRefill = tb.refillAmount
	}

	// Refill the tokens
	tb.tokens += tokensToRefill
	tb.lastRefill = tb.lastRefill.Add(time.Duration(tokensToRefill) * tb.refillRate)

	// Wait until enough tokens are available
	for tb.tokens < n {
		time.Sleep(tb.refillRate)
		tb.tokens++
	}

	// Consume the tokens
	tb.tokens -= n
}
