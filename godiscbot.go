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

var (
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

	dg.AddHandlerOnce(ready)

	addMessageHandlers(dg, BotHandlers)
	dg.Identify.Intents = botIntents

	err = dg.Open()
	if err != nil {
		log.Fatalf("failure to connect: %v", err)
	}

	return dg
}

func addMessageHandlers(dg *discordgo.Session, handlers []MessageHandler) {
	for _, h := range handlers {
		dg.AddHandler(skipSelf(h))
	}
}

func skipSelf(fn MessageHandler) MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author != s.State.User {
			fn(s, m)
		}
	}
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	botName = m.User.Username
	log.Printf("ready: using name %s", botName)
}
