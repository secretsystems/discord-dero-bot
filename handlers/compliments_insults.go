package handlers

import (
	"math/rand"
	"strings"
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
	"Can I call you the 8th wonder?",
	"I bet you make babies smile.",
	"Based",
}

var insults = []string{
	"You're a disappointment.",
	"You're about as useful as a screen door on a submarine.",
	"If your brain was dynamite, there wouldn't be enough to blow your hat off.",
	"You must have been born on a highway because that's where most accidents happen.",
	"Is your name Wi-Fi? Because I'm feeling a connection... to someone else.",
	"Roses are red, violets are blue, I have five fingers, and the middle one's for you.",
	"I'd agree with you but then we'd both be wrong.",
	"You're not the sharpest tool in the shed.",
	"You must be a parking ticket because you have 'FINE' written all over you.",
	"If you were a vegetable, you'd be a 'cute-cumber'.",
	"Are you made of copper and tellurium? Because you're Cu-Te.",
	"You're not stupid; you just have bad luck when you think.",
	"You’re about as sharp as a bowling ball.",
	"You giant piece of milk.",
	"Your face looks like it caught on fire and someone tried to put it out with a fork.",
	"If you were any more inbred, you'd be a sandwich.",
	"I'm not saying I hate you, but I would unplug your life support to charge my phone.",
	"You are definitely not the sharpest knife in the drawer.",
}

func HandleMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore messages from bots
	if message.Author.Bot {
		return
	}

	// Get the content of the message
	content := message.Content

	// Check for the !compliment command
	if strings.HasPrefix(content, "!compliment") {
		handleCompliment(discord, message)
	} else if strings.HasPrefix(content, "!insult") {
		handleInsult(discord, message)
	}
}

func getRandomPhrase(phrases []string) string {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)
	randomIndex := randomGenerator.Intn(len(phrases))
	return phrases[randomIndex]
}

func handleCompliment(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Get a random compliment
	compliment := getRandomPhrase(compliments)

	// Get mentioned users and roles
	mentionedUsers := message.Mentions
	mentionedRoles := message.MentionRoles

	if len(mentionedUsers) == 0 && len(mentionedRoles) == 0 {
		// If no users or roles are mentioned, compliment the message author
		response := message.Author.Mention() + ", " + compliment
		discord.ChannelMessageSend(message.ChannelID, response)
		return
	}

	// Create a slice to store mentioned members' and roles' mentions
	mentions := make([]string, 0)

	// Mention mentioned users
	for _, user := range mentionedUsers {
		mentions = append(mentions, user.Mention())
	}

	// Mention mentioned roles
	for _, roleID := range mentionedRoles {
		roleMention := "<@&" + roleID + ">"
		mentions = append(mentions, roleMention)
	}

	// Send the compliment to all mentioned members and roles
	response := strings.Join(mentions, " ") + ", " + compliment
	discord.ChannelMessageSend(message.ChannelID, response)
}

func handleInsult(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// Get a random insult
	insult := getRandomPhrase(insults)

	// Get mentioned user (if any)
	mentionedUsers := message.Mentions
	if len(mentionedUsers) == 0 {
		discord.ChannelMessageSend(message.ChannelID, "Are you trying to insult yourself?")
		return
	}

	// Send the insult
	response := mentionedUsers[0].Mention() + ", " + insult
	discord.ChannelMessageSend(message.ChannelID, response)
}
