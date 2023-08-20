package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var compliments = []string{
	"You're amazing!",
	"You're the best!",
	"You bring joy to everyone around you!",
	"You have a heart of gold!",
	"You're a superstar!",
	"You make the world a better place!",
	"Your positivity is infectious!",
	"You're a true inspiration!",
}

func HandleCompliment(discord *discordgo.Session, message *discordgo.MessageCreate) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	randomIndex := randomGenerator.Intn(len(compliments))
	compliment := compliments[randomIndex]

	response := fmt.Sprintf("%s", compliment)
	discord.ChannelMessageSend(message.ChannelID, response)
}
