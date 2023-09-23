package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type VoiceState struct {
	GuildID        string
	VoiceChannelID string
}

// Define errors for better error handling
var (
	ErrNoAudioURL       = errors.New("no audio URL provided")
	ErrInvalidURL       = errors.New("invalid URL")
	ErrUserNotInVC      = errors.New("user is not in a voice channel")
	ErrJoiningVC        = errors.New("error joining voice channel")
	ErrEncoding         = errors.New("error encoding audio")
	ErrFetchingPlaylist = errors.New("error fetching playlist")
	ErrReadingPlaylist  = errors.New("error reading playlist")
)

func HandleMusic(session *discordgo.Session, message *discordgo.MessageCreate) {
	command, parameters := extractCommandAndParameters(message.Content)

	switch command {
	case "!music":
		handleMusicCommand(session, message, parameters)
	}
}

func extractCommandAndParameters(content string) (string, []string) {
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func handleMusicCommand(session *discordgo.Session, message *discordgo.MessageCreate, parameters []string) {
	if len(parameters) == 0 {
		session.ChannelMessageSend(message.ChannelID, "Please provide a URL to play music.")
		return
	}

	audioURL := parameters[0]

	if audioURL == "stop" {
		stopMusic(session, message.GuildID)
		return
	}
	// Check if music is already playing in the guild
	if isMusicPlaying(session, message.GuildID) {
		session.ChannelMessageSend(message.ChannelID, "Music is already playing. Use `!music stop` to stop it.")
		return
	}
	err := playAudioFromURL(session, message.Author.ID, message.GuildID, audioURL)
	if err != nil {
		fmt.Printf("Error playing audio: %v\n", err)
	}
}

func isMusicPlaying(session *discordgo.Session, guildID string) bool {
	for _, vs := range session.VoiceConnections {
		if vs.GuildID == guildID {
			return true
		}
	}
	return false
}
func stopMusic(session *discordgo.Session, guildID string) {
	for _, vs := range session.VoiceConnections {
		if vs.GuildID == guildID {
			vs.Disconnect()
			break
		}
	}
}

func playAudioFromURL(session *discordgo.Session, userID, guildID, audioURL string) error {
	fmt.Printf("Playing audio from URL: %s\n", audioURL)

	// Validate the audio URL
	if audioURL == "" {
		return ErrNoAudioURL
	}

	parsedURL, err := url.Parse(audioURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme == "" {
		return ErrInvalidURL
	}

	// Get the user's current voice channel
	userVoiceState, err := getUserVoiceState(session, userID)
	if err != nil {
		return err
	}

	// Play the audio based on the URL type
	if strings.HasSuffix(audioURL, ".m3u") || strings.HasSuffix(audioURL, ".pls") {
		return playPlaylist(session, guildID, userVoiceState.VoiceChannelID, audioURL)
	} else if strings.Contains(audioURL, "youtu") {
		// Handle YouTube URLs
		audioURL, err := getYoutubeAudioURL(audioURL)
		if err != nil {
			return err
		}
		return playAudio(session, guildID, userVoiceState.VoiceChannelID, audioURL)
	}

	return playAudio(session, guildID, userVoiceState.VoiceChannelID, audioURL)
}

func getYoutubeAudioURL(youtubeURL string) (string, error) {
	cmd := exec.Command("yt-dlp", "-g", "-f", "bestaudio", youtubeURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func playPlaylist(session *discordgo.Session, guildID, voiceChannelID, playlistURL string) error {
	var playlistParser func(string) ([]string, error)

	if strings.HasSuffix(playlistURL, ".m3u") {
		playlistParser = parseM3UPlaylist
	} else if strings.HasSuffix(playlistURL, ".pls") {
		playlistParser = parsePlaylist
	} else {
		return ErrInvalidURL
	}

	playlistItems, err := playlistParser(playlistURL)
	if err != nil {
		return err
	}

	for _, item := range playlistItems {
		err := playAudio(session, guildID, voiceChannelID, item)
		if err != nil {
			fmt.Printf("Error playing audio from playlist: %v\n", err)
		}
	}

	return nil
}

func getUserVoiceState(session *discordgo.Session, userID string) (*VoiceState, error) {
	guilds := session.State.Guilds
	for _, guild := range guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userID {
				return &VoiceState{
					GuildID:        guild.ID,
					VoiceChannelID: vs.ChannelID,
				}, nil
			}
		}
	}
	return nil, ErrUserNotInVC
}

func playAudio(session *discordgo.Session, guildID, voiceChannelID, audioURL string) error {
	// Join the voice channel
	vc, err := session.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
	if err != nil {
		return ErrJoiningVC
	}
	defer vc.Disconnect()

	// Set up the encoding options
	opts := dca.StdEncodeOptions
	opts.RawOutput = true

	// Encode and stream the audio
	encodeSession, err := dca.EncodeFile(audioURL, opts)
	if err != nil {
		return ErrEncoding
	}
	defer encodeSession.Cleanup()

	// Start speaking
	vc.Speaking(true)

	// Stream audio frames
	for {
		frame, err := encodeSession.OpusFrame()
		if err != nil {
			break
		}

		vc.OpusSend <- frame
	}

	// Stop speaking
	vc.Speaking(false)
	return nil
}

func parsePlaylist(playlistURL string) ([]string, error) {
	resp, err := http.Get(playlistURL)
	if err != nil {
		return nil, ErrFetchingPlaylist
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrReadingPlaylist
	}

	lines := strings.Split(string(body), "\n")
	var audioURLs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "File1=") {
			audioURL := strings.TrimPrefix(line, "File1=")
			audioURLs = append(audioURLs, audioURL)
		}
	}

	return audioURLs, nil
}

func parseM3UPlaylist(playlistURL string) ([]string, error) {
	resp, err := http.Get(playlistURL)
	if err != nil {
		return nil, ErrFetchingPlaylist
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var audioURLs []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			audioURLs = append(audioURLs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, ErrReadingPlaylist
	}

	return audioURLs, nil
}
