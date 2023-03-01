package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	seenList   map[string]time.Time
	seenListMx sync.Mutex
	sentList   []*discordgo.Message

	botHandlers = []MessageHandler{
		pong,
		seeing,
		hasSeen,
		sent,
	}
)

func pong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "PONG")
	} else if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "PING")
	}
}

// TODO: something like a list. Maybe confined to an allow list of channels in .env or similar.
func seeing(s *discordgo.Session, m *discordgo.MessageCreate) {
	seenListMx.Lock()
	defer seenListMx.Unlock()
	if _, ok := seenList[m.Author.Username]; !ok {
		seenList[m.Author.Username] = m.Timestamp
	}
}

func hasSeen(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!seen") {
		return
	}

	seenListMx.Lock()
	defer seenListMx.Unlock()

	out := seenResponse(m)
	if out == "" {
		return
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		log.Printf("error sending message: %v\ncontent was: %s", err, out)
		return
	}
	sentList = append(sentList, msg)
}

func seenResponse(m *discordgo.MessageCreate) string {
	var out string
	words := strings.Fields(m.Content)
	if len(words) == 1 {
		out = fmt.Sprintf("%v", seenList)
	} else if len(words) == 2 {
		username := words[1]
		if t, ok := seenList[username]; ok {
			out = fmt.Sprintf("%s, I saw %s and it was %v", m.Message.Author, username, t)
		} else {
			out = fmt.Sprintf("%s, I've never seen %s", m.Message.Author, username)
		}
	}
	return out
}

func sent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!sent") {
		return
	}

	if len(sentList) == 0 {
		s.ChannelMessageSend(m.ChannelID, "i ... haven't seen anyone")
		return
	}

	var buf strings.Builder

	buf.WriteString("I sent these:\n")
	for i, msg := range sentList {
		out := fmt.Sprintf("%d %s %s %s\n", i, msg.ChannelID, msg.Timestamp, msg.Content)
		buf.WriteString(out)
	}

	s.ChannelMessageSend(m.ChannelID, buf.String())
}
