package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func PlayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		ChanID  = m.ChannelID
		GuildID = m.GuildID
	)

	flag.Parse()

	vs, err := s.State.VoiceState(GuildID, m.Author.ID)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You must be connected to a voice channel to play a song")
		return
	}

	vc, err := s.ChannelVoiceJoin(GuildID, ChanID, false, true)

	if err != nil {
		fmt.Println("Error: Joining the voice channel:", err)
		return
	}

	vs.Mute = false

	vc.Speaking(true)
	defer vc.Speaking(false)

	vc.Ready = true

	file, err := os.Open("./songs/currentSong.opus")
	if err != nil {
		fmt.Println("Error: Opening .opus file", err)
		return
	}

	fmt.Println("file opened")
	defer file.Close()

	buffer := make([]byte, 2048)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			fmt.Println("Error reading the file:", err)
			break
		}

		if n == 0 {
			break
		}

		vc.OpusSend <- buffer[:n]
		fmt.Println("chunk sent", n)
	}

	vc.Close()
}
