package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/thegroobi/discordgo-bot/bot/commands"
	"github.com/thegroobi/discordgo-bot/bot/helper"
)

func Start() (bot *discordgo.Session, err error) {
	godotenv.Load()
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating the discord session:", err)
		return nil, err
	}

	bot.AddHandler(messageHandler)

	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = bot.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return nil, err
	}

	fmt.Println("Starting bot...")
	return bot, nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//bot is the author of the message
	if m.Author.ID == s.State.User.ID {
		return
	}

	_, fullCommand, found := strings.Cut(m.Content, "$")

	//command doesn't start with prefix
	if !found {
		return
	}

	args := strings.Split(fullCommand, " ")
	command := strings.ToLower(args[0])

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
				fmt.Println("Error fetching pokemon:", err)
				return
			}

			if res.StatusCode == 404 {
				s.ChannelMessageSend(m.ChannelID, "Pokemon not found...")
				return
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading the response:", err)
				return
			}

			var data map[string]interface{}
			json.Unmarshal([]byte(b), &data)

			name := strings.ToUpper(data["name"].(string)[:1]) + data["name"].(string)[1:]
			pokeID := strconv.FormatFloat(data["id"].(float64), 'f', -1, 64)

			//shiny handling
			if len(args) == 3 && strings.ToLower(args[2]) == "shiny" {
				spriteShiny := data["sprites"].(map[string]interface{})["front_shiny"].(string)
				resSprite, err := http.Get(spriteShiny)
				if err != nil {
					fmt.Println("Error reading shiny GET request:", err)
					return
				}

				r := resSprite.Body

				color, err := helper.FindDominantColor(r)
				if err != nil {
					fmt.Println("Error finding the dominant color:", err)
					return
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
					fmt.Println("Error sending embed message:", err)
					return
				}
			} else {
				spriteDefault := data["sprites"].(map[string]interface{})["front_default"].(string)
				resSprite, err := http.Get(spriteDefault)
				if err != nil {
					fmt.Println("Error reading default GET request", err)
					return
				}

				r := resSprite.Body

				color, err := helper.FindDominantColor(r)
				if err != nil {
					fmt.Println("Error finding the dominant color:", err)
					return
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
					fmt.Println("Error: Sending embed message", err)
					return
				}
			}

			fmt.Printf("client: status code: %d\n", res.StatusCode)
		}

		if command == "play" {
			if len(args) < 2 {
				s.ChannelMessageSend(m.ChannelID, "No song provided.")

			} else if len(args) == 2 {

				message, err := helper.DownloadSong(args[1])

				if err != nil {
					s.ChannelMessageSend(m.ChannelID, message)
					fmt.Println(err)
					return
				}
				s.ChannelMessageSend(m.ChannelID, message)

				commands.PlayHandler(s, m)

			}
		}
	}
}