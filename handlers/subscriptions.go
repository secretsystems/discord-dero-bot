package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/secretsystems/discord-dero-bot/exports"
	"github.com/secretsystems/discord-dero-bot/utils/dero"
)

var (
	botChannelID              = "1060312629505167362"
	secretStudentsRole        = "1151545331717255199"
	privateIslandSCID         = "ce99dae86c4172378e53be91b4bb2d99f057c1eb24400510621af6002b2b10e3"
	privateIslandSubKeysuffix = "_cf530bd98d200171a94bcd6ef1e3ad6348bfa3e6691196e64e93e7953b64a2e40_E"
)

func Subscriptions(session *discordgo.Session) {
	loadUserMappings()

	now := float64(time.Now().Unix())
	for userID, userAddress := range userMappings {
		key := userAddress + privateIslandSubKeysuffix
		result := dero.GetStringKey(privateIslandSCID, key)

		switch result.(type) {
		case float64: // Fix the case type assertion
			log.Printf("User %s is subscribed", userID)

			if hasRole(session, userID) {
				if now > result.(float64) {
					removeSubscriptionsRole(session, userID)
					session.ChannelMessageSend(botChannelID, "Subscription Expired for <@"+userID+">")
					log.Printf("User %s had a role, but their subscription has expired", userID)
				} else {
					log.Printf("User %s already has a role, and their subscription is still valid", userID)
				}
			} else {
				if now < result.(float64) {
					addSubscriptionRole(session, userID)
					session.ChannelMessageSend(botChannelID, "Subscription started for <@"+userID+">")
					log.Printf("User %s didn't have a role, and their subscription is still valid", userID)
				}
			}

		default:
			// log.Printf("User %s is not subscribed", userID)
		}
	}
}

func removeSubscriptionsRole(session *discordgo.Session, userID string) {
	// Remove the registered role and add the unregistered role
	err := session.GuildMemberRoleRemove(exports.SecretGuildID, userID, secretStudentsRole)
	if err != nil {
		log.Println("Error removing role from member:", err)
	}
}

func addSubscriptionRole(session *discordgo.Session, userID string) {
	err := session.GuildMemberRoleAdd(exports.SecretGuildID, userID, secretStudentsRole)
	if err != nil {
		log.Printf("Error adding role for Guild %v to member:%v", exports.SecretMembersRoleID, err)
	}
}

func hasRole(session *discordgo.Session, userID string) bool {
	member, err := session.GuildMember(exports.SecretGuildID, userID)
	if err != nil {
		log.Printf("Error getting guild member: %v", err)
		return false
	}
	for _, roleID := range member.Roles {
		switch roleID {
		case secretStudentsRole:
			return true
		}
	}
	return false
}

func HandlePrivateSubscriptions(session *discordgo.Session, message *discordgo.MessageCreate) {
	loadUserMappings()
	content := message.Content

	if !strings.HasPrefix(content, "!dev101 ") {
		userMessage := "To lookup DeroDev101 Subscription, use the format: `!dev101 <dero1qaddress>`"
		session.ChannelMessageSend(message.ChannelID, userMessage)
		return
	}

	userInput := strings.TrimPrefix(message.Content, "!dev101 ")

	key := userInput + privateIslandSubKeysuffix
	result := dero.GetStringKey(privateIslandSCID, key).(float64)

	// Convert Unix timestamp to human-readable format
	timestamp := time.Unix(int64(result), 0)
	formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

	// Send the result to the channel
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Subscription Expires on %s", formattedTimestamp))
}
