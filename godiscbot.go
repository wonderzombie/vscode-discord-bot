package main

import (
	"log"

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
	session := start(tok)
	defer session.Close()

	bot.New(session).Run()
}

func start(tok string) *discordgo.Session {
	session, err := discordgo.New("Bot " + tok)
	if err != nil {
		log.Fatalf("failure to launch: %v", err)
	}

	session.Identify.Intents = botIntents

	err = session.Open()
	if err != nil {
		log.Fatalf("failure to connect: %v", err)
	}

	return session
}
