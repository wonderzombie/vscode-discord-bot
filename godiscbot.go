package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wonderzombie/godiscbot/bot"
)

var (
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
	defer dg.Close()

	fmt.Println("running...")
	awaitTerm()
}

func awaitTerm() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func start(tok string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + tok)
	if err != nil {
		log.Fatalf("failure to launch: %v", err)
	}

	bot.AddMessageHandlers(dg)
	dg.Identify.Intents = botIntents

	err = dg.Open()
	if err != nil {
		log.Fatalf("failure to connect: %v", err)
	}
	return dg
}
