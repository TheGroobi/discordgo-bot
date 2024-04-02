package helper

import (
	"bytes"
	"fmt"
	"os/exec"
)

func DownloadSong(url string) (message string, err error) {
	var stdout, stderr bytes.Buffer

	c := exec.Command("cmd.exe", "/c", "/Design/discordgo-bot/bot/helper/download-song/ffmpeg.sh", url)
	fmt.Println("Executing bash script...")

	c.Stdout = &stdout
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		fmt.Println("Error: Running bash", err)
		return "Something went wrong...", err
	}
	fmt.Println("Song encoded correctly")

	return "Song downloaded successfully!", err
}
