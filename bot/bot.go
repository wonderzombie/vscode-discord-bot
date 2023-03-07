package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

////// CORE BOT FUNCTIONALITY /////

// DiscordBot is a glorified container for a discordgo Session.
type DiscordBot struct {
	s    *discordgo.Session
	User *discordgo.User
	mods []*Module
	sent []*discordgo.Message
}

func New(s *discordgo.Session, modules ...Handler) *DiscordBot {
	b := &DiscordBot{
		s:    s,
		User: s.State.User,
	}

	s.AddHandlerOnce(Ready)
	s.AddHandler(b.main)

	return b
}

// Run is meant to block the main thread while the discordgo API manages handlers.
func (b *DiscordBot) Run() {
	fmt.Println("running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("quitting")
}

func (b *DiscordBot) main(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.String() == m.Author.String() {
		return
	}

	msg := NewMessage(m)
	for _, module := range b.mods {
		if fired, resp := module.Handle(msg); fired {
			for _, o := range resp {
				b.sendMessage(m.ChannelID, o)
			}
		}
	}
}

func (b *DiscordBot) sendMessage(channelID string, o string) {
	msg, err := b.s.ChannelMessageSend(channelID, o)
	if err != nil {
		fmt.Println("unable to send message:", err, "\nMessage follows: [", o, "]")
	} else {
		b.sent = append(b.sent, msg)
	}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("ready: using name %s", s.State.User.Username)
}
