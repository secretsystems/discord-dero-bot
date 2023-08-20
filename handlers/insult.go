package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var insults = []string{
	"You're not the sharpest tool in the shed.",
	"You must be a parking ticket because you have 'FINE' written all over you.",
	"If you were a vegetable, you'd be a 'cute-cumber'.",
	"Are you made of copper and tellurium? Because you're Cu-Te.",
	"You're not stupid; you just have bad luck when you think.",
	"Your face looks like it caught on fire and someone tried to put it out with a fork.",
	"If you were any more inbred, you'd be a sandwich.",
	"I'm not saying I hate you, but I would unplug your life support to charge my phone.",
}

func HandleInsult(discord *discordgo.Session, message *discordgo.MessageCreate) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	randomIndex := randomGenerator.Intn(len(insults))
	insult := insults[randomIndex]

	response := fmt.Sprintf("%s", insult)
	discord.ChannelMessageSend(message.ChannelID, response)
}
