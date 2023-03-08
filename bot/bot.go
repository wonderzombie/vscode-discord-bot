package bot

import (
	"fmt"
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

type Bot interface {
	Send([]string) []*discordgo.Message
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
	fmt.Println("[bot] running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("[bot] quitting")
}

type SendStatus struct {
	message *discordgo.Message
	err     error
}

type StatusAll []*SendStatus

func (sa StatusAll) ok() bool {
	for _, ss := range sa {
		if ss.err != nil {
			return false
		}
	}
	return true
}

func (b *DiscordBot) Send(channelID string, lines []string) {
	allStatus := make(StatusAll, len(lines))
	for _, l := range lines {
		allStatus = append(allStatus, b.sendMessage(channelID, l))
	}
	fmt.Println("[bot] all send ok?", allStatus.ok())
}

func (b *DiscordBot) sendMessage(channelID string, o string) *SendStatus {
	msg, err := b.s.ChannelMessageSend(channelID, o)
	if err != nil {
		fmt.Println("[bot] unable to send message:", err, "\nMessage follows: [", o, "]")
	} else {
		b.sent = append(b.sent, msg)
	}
	return &SendStatus{msg, err}
}

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	fmt.Printf("[bot] ready: using name %s\n", s.State.User.Username)
}

func (b *DiscordBot) messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.runResponders(NewMessage(m))
}

func (b *DiscordBot) runResponders(m *Message) {
	if m.Author == b.User.String() {
		return
	}

	var numFired int
	for _, r := range b.responders {
		if fired, resp := r(m); fired && !Empty(resp) {
			numFired++
			for _, out := range resp {
				status := b.sendMessage(m.ChannelID, out)
				if status.err != nil {
					fmt.Printf("[bot] error sending message: %v", status.err)
				}
			}
		}
	}
	if numFired > 0 {
		fmt.Printf("[bot] handlers fired %d", numFired)
	}
}
