package handlers

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/bwmarrin/discordgo"
// )

// var compliments = []string{

// }

// func HandleCompliment(discord *discordgo.Session, message *discordgo.MessageCreate) {
// 	randomSource := rand.NewSource(time.Now().UnixNano())
// 	randomGenerator := rand.New(randomSource)

// 	randomIndex := randomGenerator.Intn(len(compliments))
// 	compliment := compliments[randomIndex]

// 	response := fmt.Sprintf("%s", compliment)
// 	discord.ChannelMessageSend(message.ChannelID, response)
// }
