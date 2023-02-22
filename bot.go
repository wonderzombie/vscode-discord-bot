package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler func(*discordgo.Session, *discordgo.MessageCreate)

var seenList = make(map[string]time.Time)

var sent = make([]*discordgo.Message, 20)

var seenListMx = sync.Mutex{}

func Seeing(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := seenList[m.Author.Username]; !ok {
		seenListMx.Lock()
		defer seenListMx.Unlock()
		seenList[m.Author.Username] = m.Timestamp
	}
}

func HasSeen(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!seen") {
		return
	}

	var out string
	words := strings.Fields(m.Content)
	if len(words) == 1 {
		out = fmt.Sprintf("%v", seenList)
	} else if len(words) == 2 {
		username := words[1]
		seenListMx.Lock()
		defer seenListMx.Unlock()
		if t, ok := seenList[username]; ok {
			out = fmt.Sprintf("%s, I saw %s and it was %v", m.Message.Author, username, t)
		} else {
			out = fmt.Sprintf("%s, I've never seen %s", m.Message.Author, username)
		}
	}

	if out == "" {
		return
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		log.Printf("error sending message: %v\ncontent was: %s", err, out)
		return
	}
	sent = append(sent, msg)
}

func Sent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, "!sent") {
		return
	}

	if len(sent) == 0 {
		s.ChannelMessageSend(m.ChannelID, "i ... haven't seen anyone")
		return
	}

	var buf strings.Builder

	buf.WriteString("I sent these:\n")

	for i, msg := range sent {
		out := fmt.Sprintf("%d %s %s %s\n", i, msg.ChannelID, msg.Timestamp, msg.Content)
		buf.WriteString(out)
	}

	s.ChannelMessageSend(m.ChannelID, buf.String())
}
