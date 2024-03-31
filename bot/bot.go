package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func Start() *discordgo.Session {
	godotenv.Load()
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	bot.AddHandler(messageHandler)

	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = bot.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is running...")

	return bot
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	_, fullCommand, found := strings.Cut(m.Content, "$")

	if !found {
		return
	}

	args := strings.Split(strings.ToLower(fullCommand), " ")
	command := args[0]

	//hello world command
	if command == "hello" {
		s.ChannelMessageSend(m.ChannelID, "World!")
	}

	if command == "miki" {
		s.ChannelMessageSend(m.ChannelID, "Mikołaj giga fiut")
	}

	if len(args) >= 2 {
		//PokeAPI fetch command
		//send embed picture with name of fetched pokemon
		if command == "poke" {
			res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + args[1])

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Something went wrong...")
				log.Fatal(err)
			}

			if res.StatusCode == 404 {
				s.ChannelMessageSend(m.ChannelID, "Pokemon not found...")
				return
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			var data map[string]interface{}
			json.Unmarshal([]byte(b), &data)

			name := strings.ToUpper(data["name"].(string)[:1]) + data["name"].(string)[1:]
			pokeID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)

			//shiny handling
			if len(args) == 3 && args[2] == "shiny" {
				spriteShiny := data["sprites"].(map[string]interface{})["front_shiny"].(string)
				resSprite, err := http.Get(spriteShiny)
				if err != nil {
					log.Fatal(err)
				}

				r := resSprite.Body

				color, err := FindDominantColor(r)
				if err != nil {
					log.Fatal(err)
				}

				shinyEmbed := &discordgo.MessageEmbed{
					Title:       name + " " + strings.ToUpper(args[2][:1]) + args[2][1:],
					Description: "Pokédex ID: " + pokeID,
					Color:       int(color),
					Image: &discordgo.MessageEmbedImage{
						URL: spriteShiny,
					},
				}

				_, err = s.ChannelMessageSendEmbeds(m.ChannelID, []*discordgo.MessageEmbed{shinyEmbed})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				spriteDefault := data["sprites"].(map[string]interface{})["front_default"].(string)
				resSprite, err := http.Get(spriteDefault)
				if err != nil {
					log.Fatal(err)
				}

				r := resSprite.Body

				color, err := FindDominantColor(r)
				if err != nil {
					log.Fatal(err)
				}

				defaultEmbed := &discordgo.MessageEmbed{
					Title:       name,
					Description: "Pokédex ID: " + pokeID,
					Color:       int(color),
					Image: &discordgo.MessageEmbedImage{
						URL: spriteDefault,
					},
				}

				_, err = s.ChannelMessageSendEmbeds(m.ChannelID, []*discordgo.MessageEmbed{defaultEmbed})
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Printf("client: status code: %d\n", res.StatusCode)
		}

		if command == "play" {
			if len(args) == 1 {
				s.ChannelMessageSend(m.ChannelID, "No song provided.")
			} else if len(args) == 2 {
				downloadSong(args[1])
			}
		}
	}
}

func downloadSong(url string) {

	println("executing bash")
	c := exec.Command("./downloadSong.sh", url)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}

}
