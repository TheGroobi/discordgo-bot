package commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thegroobi/discordgo-bot/bot/helper"
)

var buffer = make([][]byte, 0)

func PlayHandler(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {

	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		helper.OnError("Establishing voice state", err)
		s.ChannelMessageSend(m.ChannelID, "You must be connected to a voice channel to play a song")
		return err
	}

	loadSong()
	if err != nil {
		helper.OnError("Loading song", err)
		return err
	}

	err = playSong(s, m.GuildID, vs.ChannelID)
	if err != nil {
		helper.OnError("Playing song", err)
	}

	return nil
}

func loadSong() error {
	file, err := os.Open("songs/currentSong.opus")
	if err != nil {
		helper.OnError("Opening dca file", err)
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			helper.OnError("Closing the file", err)
		}
	}()

	// var opuslen int8
	counter := 0
	for {

		counter++
		fmt.Printf("file read 1, times:%d\n", counter)
		// err = binary.Read(file, binary.LittleEndian, &opuslen)
		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	helper.OnError("Reading from dca file", err)
		// 	return err
		// }

		InBuf := make([]byte, 960)
		err := binary.Read(file, binary.LittleEndian, &InBuf)

		fmt.Printf("file read 2, times:%d\n", counter)
		if err != nil {
			if err == io.EOF {
				break
			}
			helper.OnError("Reading from dca file", err)
			return err
		}

		buffer = append(buffer, InBuf)
	}
	return nil
}

func playSong(s *discordgo.Session, gID, cID string) error {

	vc, err := s.ChannelVoiceJoin(gID, cID, false, true)
	if err != nil {
		helper.OnError("Connecting to voice", err)
		return err
	}

	time.Sleep(250 * time.Millisecond)

	err = vc.Speaking(true)
	if err != nil {
		helper.OnError("Set speaking", err)
		return err
	}

	time.Sleep(250 * time.Millisecond)

	defer func() {
		err = vc.Speaking(false)
		if err != nil {
			helper.OnError("Couldn't set speaking", err)
			return
		}
		vc.Disconnect()
	}()

	for _, b := range buffer {
		vc.OpusSend <- b
	}
	return nil
}
