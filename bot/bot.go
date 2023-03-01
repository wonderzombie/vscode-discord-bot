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
	Session Sesh
}

type Sesh interface {
	AddHandler(interface{}) func()
	AddHandlerOnce(interface{}) func()
}

func New(s Sesh) *Bot {
	b := &Bot{
		Session: s,
	}
	initHandlers(s)
	return b
}

func (b *Bot) Run() {
	fmt.Println("running")
	b.awaitTerm()
	fmt.Println("quitting")
}

func initHandlers(s Sesh) {
	for _, h := range botHandlers {
		s.AddHandler(skipSelf(h))
	}
	s.AddHandlerOnce(Ready)
}

func (b *Bot) awaitTerm() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	botName = m.User.Username
	log.Printf("ready: using name %s", botName)
}

func skipSelf(mh MessageHandler) MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author != s.State.User {
			mh(s, m)
		}
	}
}
