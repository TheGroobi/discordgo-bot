package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/thegroobi/discordgo-bot/bot"
)

func main() {
	bot := bot.Start()

	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
