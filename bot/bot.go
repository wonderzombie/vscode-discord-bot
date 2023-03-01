package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// MessageHandler describes a single function which can act as a handler in discordgo.
type MessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

var (
	botName = ""
)

type Bot struct {
	Session *discordgo.Session
}

func New(dg *discordgo.Session) *Bot {
	b := &Bot{
		Session: dg,
	}

	b.addHandlers(botHandlers)
	b.Session.AddHandlerOnce(Ready)

	return b
}

func (b *Bot) Run() {
	fmt.Println("running")
	b.awaitTerm()
	fmt.Println("quitting")
}

func (b *Bot) addHandlers(msgHandlers []MessageHandler) {
	for _, h := range botHandlers {
		b.Session.AddHandler(h)
	}
	for _, mh := range msgHandlers {
		b.Session.AddHandler(skipSelf(mh))
	}
}

func (b *Bot) awaitTerm() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func AddMessageHandlers(dg *discordgo.Session) {
	dg.AddHandlerOnce(Ready)
	for _, h := range botHandlers {
		dg.AddHandler(skipSelf(h))
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
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
