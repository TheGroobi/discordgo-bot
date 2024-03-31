package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thegroobi/discordgo-bot/bot"
)

func main() {
	bot, err := bot.Start()

	if err != nil {
		fmt.Println("Error starting the bot:", err)
		return
	}

	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	
}
