package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

///// MESSAGE HANDLING /////

type Responder func(m *Message) (bool, []string)

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
	return nil, false
}

func NewMessage(m *discordgo.MessageCreate) *Message {
	fields := strings.Fields(m.Content)
	return &Message{m.Author.String(), m.ChannelID, m.Content, fields}

}
