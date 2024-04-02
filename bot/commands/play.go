package commands

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thegroobi/discordgo-bot/bot/helper"
)

func PlayHandler(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {

	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		helper.OnError("Establishing voice state", err)
		s.ChannelMessageSend(m.ChannelID, "You must be connected to a voice channel to play a song")
		return err
	}

	err = playSong(s, m.GuildID, vs.ChannelID)
	if err != nil {
		helper.OnError("Playing song", err)
	}

	return nil
}

// func loadSong(vc *discordgo.VoiceConnection) error {
// 	file, err := os.Open("/test.dca")
// 	if err != nil {
// 		helper.OnError("Opening dca file", err)
// 		return err
// 	}

// 	defer func() {
// 		err := file.Close()
// 		if err != nil {
// 			helper.OnError("Closing the file", err)
// 		}
// 	}()

// 	// var opuslen int16
// 	counter := 0

// 	for {

// 		counter++
// 		fmt.Printf("file read 1, times:%d\n", counter)

// 		// err = binary.Read(file, binary.LittleEndian, &opuslen)
// 		chunk, err := file.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			helper.OnError("Reading from dca file", err)
// 			return err
// 		}

// 		vc.OpusSend <- buffer[:chunk]
// 		// 	err := binary.Read(file, binary.LittleEndian, &InBuf)

// 		// 	if err != nil {
// 		// 		if err == io.EOF {
// 		// 			break
// 		// 		}
// 		// 		helper.OnError("Reading from dca file", err)
// 		// 		return err
// 		// 	}

// 		// 	buffer = append(buffer, InBuf)
// 		// }
// 	}
// 	return nil
// }

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

	// defer func() {
	// 	err = vc.Speaking(false)
	// 	if err != nil {
	// 		helper.OnError("Couldn't set speaking", err)
	// 		return
	// 	}
	// 	vc.Disconnect()
	// }()

	file, err := os.Open("bot/helper/download-song/songs/output.opus")
	if err != nil {
		helper.OnError("Opening opus file", err)
		return err
	}

	// defer func() {
	// 	err := file.Close()
	// 	if err != nil {
	// 		helper.OnError("Closing the file", err)
	// 	}
	// }()

	counter := 0
	var buffer = make([]byte, 960*2*2)

	for {

		counter++
		fmt.Printf("read: %d times\n", counter)

		// err = binary.Read(file, binary.LittleEndian, &opuslen)
		chunk, err := file.Read(buffer)

		// fmt.Printf("\nchunk %d: %s", chunk, buffer[:chunk])
		if err != nil {
			if err == io.EOF {
				break
			}
			helper.OnError("Reading from dca file", err)
			return err
		}

		if chunk == len(buffer) {
			vc.OpusSend <- buffer[:chunk]
		}
	}
	return nil
}

// for _, buf := range buffer {
// 	vc.OpusSend <- buf
// 	if err != nil {
// 		if err == io.EOF {
// 			break
// 		}
// 		helper.OnError("Reading from dca file", err)
// 		return err
// 	}
// }
