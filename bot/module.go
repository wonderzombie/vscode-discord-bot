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

func (m *Message) Fields() []string {
	return m.fields
}

func (m *Message) Cmd() (string, bool) {
	if f := m.fields[0]; strings.HasPrefix(f, "!") {
		return f, true
	}
	return "", false
}

func (m *Message) Args() ([]string, bool) {
	if len(m.fields) > 1 {
		out := m.fields[1:len(m.fields)]
		return out, true
	}
	return []string{""}, false
}

func NewMessage(m *discordgo.MessageCreate) *Message {
	fields := strings.Fields(m.Content)
	return &Message{m.Author.String(), m.ChannelID, m.Content, fields}

}

type Processor func(m *Message) bool
type Responder func(m *Message) (bool, []string)
