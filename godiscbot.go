package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	envVars, err := godotenv.Read()
	if err != nil {
		log.Fatalf("failure to load config: %v", err)
	}

	tok, ok := envVars["BOT_TOKEN"]
	if !ok || tok == "" {
		log.Fatalf(".env is missing a required field: BOT_TOKEN")
	}
	dg := start(tok)

	fmt.Println("running...")

	awaitClose()
	dg.Close()
}

func awaitClose() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func start(tok string) (s *discordgo.Session) {
	dg, err := discordgo.New(tok)
	if err != nil {
		log.Fatalf("failure to launch: %v", err)
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf("failure to connect: %v", err)
	}
	return s
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "PONG")
	} else if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "PING")
	}
}
