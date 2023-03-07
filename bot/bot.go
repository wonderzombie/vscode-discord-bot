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

type MessageCreateResponder func(*discordgo.MessageCreate) (bool, []string)
type MessageResponder func(*Message) (bool, []string)

///// MESSAGE HANDLING /////

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

///// MODULES /////

type Handler interface {
	Handle(m *discordgo.MessageCreate) Response
}

type Module struct {
	Id        string
	responder MessageResponder
}

func NewModule(id string, mr MessageResponder) *Module {
	return &Module{
		Id:        id,
		responder: mr,
	}
}

func (m *Module) Handle(msg *Message) (bool, []string) {
	return m.responder(msg)
}

////// BOT FUNCTIONALITY /////

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
