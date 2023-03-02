package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type SessionHandler func(*discordgo.Session, *discordgo.MessageCreate)
type ReadyHandler func(*discordgo.Session, *discordgo.Ready)

type MessageHandler func(*discordgo.MessageCreate)
type MessageResponder func(*discordgo.MessageCreate) []string

// DiscordBot is a glorified container for a discordgo Session.
type DiscordBot struct {
	s    *discordgo.Session
	User *discordgo.User
}

// Run is meant to block the main thread while the discordgo API manages handlers.
func (b *DiscordBot) Run() {
	fmt.Println("running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("quitting")
}

func New(s *discordgo.Session) *DiscordBot {
	initHandlers(s)

	return &DiscordBot{
		s:    s,
		User: s.State.User,
	}
}

func initHandlers(s *discordgo.Session) {
	s.AddHandlerOnce(Ready)

	for _, rd := range readiers {
		s.AddHandlerOnce(rd)
	}

	for _, rs := range responders {

		s.AddHandler(rs)
	}

	for _, h := range handlers {
		s.AddHandler(h)
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("ready: using name %s", s.State.User.Username)
}

func (ss *SeenState) skipSelf(username string, fn MessageResponder) MessageResponder {
	return func(m *discordgo.MessageCreate) []string {
		ret := make([]string, 0)
		if m.Author.Username != ss.user.Username {
			ret = fn(m)
		}
		return ret
	}
}

func empty(lines []string) bool {
	ret := false
	if len(lines) == 0 {
		ret = true
	} else if l := len(lines); l == 1 && lines[0] == "" {
		ret = true
	}
	return ret
}
