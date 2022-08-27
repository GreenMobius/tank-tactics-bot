package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const cmdPrefix rune = '!'

const helpText string = `Available commands:
	ping    send a ping and wait for a pong
	help    display help text`

var botToken string

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
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, helpText)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	}
}

func main() {
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + botToken)
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
