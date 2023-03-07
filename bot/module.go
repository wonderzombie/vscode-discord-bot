package bot

import (
	"strings"

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

type Processor func(m *Message) bool
type Responder func(m *Message) (bool, []string)
