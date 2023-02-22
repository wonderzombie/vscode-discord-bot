package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	names      = make(map[string]time.Time)
	botName    = ""
	botIntents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentGuildMessageReactions | discordgo.IntentDirectMessageReactions
)

func main() {
	envVars, err := godotenv.Read()
	if err != nil {
		log.Fatalf("failure to load config: %v", err)
	}

	tok, ok := envVars["DISCORD_TOKEN"]
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

func start(tok string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + tok)
	if err != nil {
		log.Fatalf("failure to launch: %v", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandlerOnce(ready)
	dg.Identify.Intents = botIntents

	err = dg.Open()
	if err != nil {
		log.Fatalf("failure to connect: %v", err)
	}

	return dg
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	botName = m.User.Username
	log.Printf("ready: using name %s", botName)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.ID == "" {
		return
	}

	// Add some logging here to check what's what
	log.Printf("messageCreate: %s - %s - %s\n", m.Timestamp, m.Author, m.Content)
	if _, ok := names[m.Author.ID]; !ok {
		names[m.Author.ID] = m.Timestamp
		fmt.Printf("added name %s to registry", m.Author.ID)
	}

	if mentioned := strings.Contains(m.Content, botName); mentioned {
		content := fmt.Sprintf("where is %s, where is %s", botName, botName)
		s.ChannelMessageSend(m.ChannelID, content)
		return
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "PONG")
	} else if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "PING")
	}
}
