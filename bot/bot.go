package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type SessionHandler func(*discordgo.Session, *discordgo.MessageCreate)
type ReadyHandler func(*discordgo.Session, *discordgo.Ready)

type MessageHandler func(*discordgo.MessageCreate)
type MessageResponder func(*discordgo.MessageCreate) []string

// DiscordBot is a glorified container for a discordgo Session.
type Message struct {
	Author    string
	ChannelID string
	Content   string
	fields    []string
}

func (m *Message) Cmd() (string, bool) {
	if f := m.fields[0]; strings.HasPrefix(f, "!") {
		return f, true
	}
	return "", false
}

func NewMessage(m *discordgo.MessageCreate) *Message {
	fields := strings.Fields(m.Content)
	return &Message{
		Author:    m.Author.String(),
		ChannelID: m.ChannelID,
		Content:   m.Content,
		fields:    fields,
	}
}

type Response struct {
	HasOut bool
	Out    []string
	Final  bool
}

type Handler interface {
	Handle(m *Message) Response
}

type BotModule struct{}

func (m *BotModule) Handle(msg *Message) Response {
	return Response{false, []string{}, false}
}

type DiscordBot struct {
	s    *discordgo.Session
	User *discordgo.User
	mods []*BotModule
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

	var out []string
	msg := NewMessage(m)
	for _, module := range b.mods {
		resp := module.Handle(msg)
		if resp.HasOut {
			out = resp.Out
		}
		if resp.Final {
			break
		}
	}

	for _, o := range out {
		b.sendMessage(m.ChannelID, o)
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

func empty(lines []string) bool {
	ret := false
	if len(lines) == 0 {
		ret = true
	} else if l := len(lines); l == 1 && lines[0] == "" {
		ret = true
	}
	return ret
}
