package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

var (
	botName = ""
)

type Bot struct {
	dg *discordgo.Session
}

func New(dg *discordgo.Session, ready MessageHandler) *Bot {
	dg.AddHandlerOnce(ready)
	return &Bot{
		dg: dg,
	}
}

func (b *Bot) AddHandlers(handlers []MessageHandler) {
	for _, h := range handlers {
		b.dg.AddHandler(skipSelf(h))
	}
}

func AddMessageHandlers(dg *discordgo.Session) {
	dg.AddHandlerOnce(ready)
	for _, h := range Handlers {
		dg.AddHandler(skipSelf(h))
	}
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	botName = m.User.Username
	log.Printf("ready: using name %s", botName)
}

func skipSelf(fn MessageHandler) MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author != s.State.User {
			fn(s, m)
		}
	}
}
