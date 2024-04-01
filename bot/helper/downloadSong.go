package helper

import (
	"bytes"
	"fmt"
	"os/exec"
)

func DownloadSong(url string) (message string, err error) {

	c := exec.Command("python", "/Design/discordgo-bot/bot/download-song/main.py", url)
	fmt.Println("Executing python script...")

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Python Error:", stderr.String())
		return "Something went wrong...", err
	}
	fmt.Println("Song downloaded")

	return "Song downloaded successfully!", err
}
