package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wonderzombie/godiscbot/bot"
	"github.com/wonderzombie/godiscbot/combat"
	"github.com/wonderzombie/godiscbot/reply"
	"github.com/wonderzombie/godiscbot/seen"
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

	bot.New(session,
		combat.Responder(),
		seen.Responder(),
		reply.Responder(session.State.User),
	).Run()

	os.Exit(0)
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
