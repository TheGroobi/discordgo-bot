package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thegroobi/discordgo-bot/bot"
	"github.com/thegroobi/discordgo-bot/bot/helper"
)

func main() {
	var err error

	bot, err := bot.Start()

	if err != nil {
		helper.OnError("Starting the bot", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	
	bot.Close()
}
