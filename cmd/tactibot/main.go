package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/greenmobius/tank-tactics-bot/pkg/tactics"
)

const cmdPrefix rune = '!'

const helpText string = `Available commands:
	ping                    send a ping and wait for a pong
	new <width> <height>    create a new map of specified width and height
	map                     display the game map
	set <x> <y> <color>     set a map location to a specified color
	help                    display help text`

var gameMap tactics.Map

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore commands not prefixed with the command prefix
	if len(m.Content) == 0 || rune(m.Content[0]) != cmdPrefix {
		return
	}

	command := strings.Split(m.Content[1:], " ")
	log.Printf("Received command: %v\n", command)
	switch command[0] {
	case "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	case "new":
		if len(command) != 3 {
			_, err := s.ChannelMessageSend(m.ChannelID, helpText)
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		width, err := strconv.ParseInt(command[1], 10, 32)
		if err != nil || width < 1 || width > 20 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Width must be a number between 1 and 20")
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		height, err := strconv.ParseInt(command[2], 10, 32)
		if err != nil || height < 1 || height > 20 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Height must be a number between 1 and 20")
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		gameMap = tactics.NewMap(int(width), int(height))
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Created new map of width %d and height %d", width, height))
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	case "map":
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Game Status",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Map",
					Value: gameMap.ToDiscordString(),
				},
			},
		})
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	case "set":
		if len(command) != 3 {
			_, err := s.ChannelMessageSend(m.ChannelID, helpText)
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		x, err := strconv.ParseInt(command[1], 10, 32)
		if err != nil || x < 0 || int(x) > len(gameMap.Grid) {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("X must be a number between 0 and %d", len(gameMap.Grid)))
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		y, err := strconv.ParseInt(command[2], 10, 32)
		if err != nil || y < 0 || int(y) > len(gameMap.Grid[0]) {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Y must be a number between 0 and %d", len(gameMap.Grid[0])))
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
			}
		}

		gameMap.Grid[x][y].Value = command[3]
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, helpText)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	}
}

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Printf("Error creating discord session: %v\n", err)
		return
	}

	dg.AddHandler(handleMessage)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Printf("Error opening discord websocket: %v\n", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot initialized...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
