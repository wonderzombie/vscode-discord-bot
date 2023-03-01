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

// Bot is a glorified container for a discordgo Session.
type Bot struct {
	Session  *discordgo.Session
	Username string
}

// Run is meant to block the main thread while the discordgo API manages handlers.
func (b *Bot) Run() {
	fmt.Println("running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("quitting")
}

func New(s *discordgo.Session) *Bot {
	initHandlers(s)
	return &Bot{
		Session:  s,
		Username: s.State.User.Username,
	}
}

func initHandlers(s *discordgo.Session) {
	s.AddHandlerOnce(Ready)
	for _, h := range botHandlers {
		s.AddHandler(skipSelf(h))
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("ready: using name %s", s.State.User.Username)
}

func skipSelf(fn MessageHandler) MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author != s.State.User {
			fn(s, m)
		}
	}
}
