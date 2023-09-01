package handlers

import (
	"discord-dero-bot/utils/dero"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Start monitoring for transfers when the user initiates membership
func startMonitoringForTransfer(session *discordgo.Session, guildID, userID string) {
	fmt.Fprintf(os.Stdout, "WE ARE STARTING SCAN FOR %s TO PURCHASE\n", userID)

	go checkForTransfer(session, guildID, userID)
}

func checkForTransfer(session *discordgo.Session, guildID, userID string) {
	fmt.Fprintf(os.Stdout, "The desiredRole: %s\n", desiredRole)

	// Define a flag to track whether a transfer was found
	transferFound := false

	// Calculate the end time for scanning (10 minutes from now)
	endTime := time.Now().Add(10 * time.Minute)

	for {
		// Perform the wallet transfer check here using the FilterAndPrintTransactions function
		fmt.Fprintf(os.Stdout, "WE ARE GOING TO FILTER FOR %s\n", userID)
		found, err := dero.FilterAndPrintTransactions(userID, deroMembershipAmount)
		fmt.Fprintf(os.Stdout, "WE HAVE SCANNED FOR %s\n", userID)
		if err != nil {
			log.Printf("Error checking for transfer: %v", err)
		} else if found {
			// Check if the user has the desired role
			member, err := session.GuildMember(guildID, userID)
			if err != nil {
				log.Printf("Error fetching member info: %v", err)
			}
			hasRole := false
			for _, roleID := range member.Roles {
				if roleID == desiredRole {
					hasRole = true
					session.ChannelMessageSend(resChan, "<@"+userID+"> is already a member!")
					break
				}
				if !hasRole {
					// Add the desired role to the member
					err := session.GuildMemberRoleAdd(guildID, userID, desiredRole)
					if err != nil {
						log.Printf("Error adding role to member: %v", err)
					} else {
						// Set the flag to true if a transfer is found
						transferFound = true
						session.ChannelMessageSend(resChan, "<@"+userID+"> is now a member!")

					}

				}

			}
		}

		// Check if the current time has passed the end time
		if time.Now().After(endTime) {
			fmt.Fprintf(os.Stdout, "SCAN TIME EXCEEDED FOR %s\n", userID)
			return
		}
		// Check if a transfer was found during the scan, if so, stop checking
		if transferFound {
			fmt.Fprintf(os.Stdout, "TRANSFER FOUND FOR %s\n", userID)
			return
		}

		// Sleep for a while before the next scan (e.g., 18 seconds)
		time.Sleep(18 * time.Second)

	}
}
