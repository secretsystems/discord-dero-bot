package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

var voiceChannelID = "1154416606361948192"

func HandleMusic(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if the message is "!music"
	if message.Content == "!music" {
		// Replace "voiceChannelID" with the actual voice channel ID

		// Join the voice channel
		vc, err := session.ChannelVoiceJoin(message.GuildID, voiceChannelID, false, true)
		if err != nil {
			fmt.Println("Error joining voice channel:", err)
			return
		}

		// Stream the audio from the URL
		audioURL := "https://ice4.somafm.com/thetrip-128-aac" // AAC 128kb link
		err = playAudioFromURL(vc, audioURL)
		if err != nil {
			fmt.Println("Error streaming audio:", err)
			return
		}
	}
}

func playAudioFromURL(vc *discordgo.VoiceConnection, audioURL string) error {
	opts := dca.StdEncodeOptions

	// Create an encoding session
	encodeSession, err := dca.EncodeFile(audioURL, opts)
	if err != nil {
		return err
	}
	defer encodeSession.Cleanup()

	vc.Speaking(true)

	// Read audio frames and send them to Discord
	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			break
		}

		vc.OpusSend <- frame
	}

	vc.Speaking(false)
	return nil
}
