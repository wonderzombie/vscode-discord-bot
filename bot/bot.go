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
	s          *discordgo.Session
	User       *discordgo.User
	responders []Responder
	sent       []*discordgo.Message
}

func New(s *discordgo.Session, responders ...Responder) *DiscordBot {
	b := &DiscordBot{
		s:          s,
		User:       s.State.User,
		responders: responders,
	}

	s.AddHandlerOnce(Ready)
	// messageCreated will invoke module handlers accordingly.
	s.AddHandler(b.messageCreated)

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

func (b *DiscordBot) messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.String() == m.Author.String() {
		return
	}

	for _, h := range b.responders {
		if fired, resp := h(NewMessage(m)); fired && !Empty(resp) {
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
