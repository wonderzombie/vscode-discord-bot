package recall

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

type recallMod struct {
	nick string
	rec  recollection
}

func (rm *recallMod) Responder(m *bot.Message) (fired bool, out []string) {
	nick := fmt.Sprintf("%s,", strings.ToLower(rm.nick))
	if !strings.HasPrefix(strings.ToLower(m.Content), nick) {
		return
	}

	// TODO: add a memory
	outStr := "I've never heard anything about that"
	if !strings.HasSuffix("?", m.Content) {
		outStr = "Good to know!"
		fired = true
		rm.rec.add(m, time.Now().Local())
	} else if mem, ok := rm.rec.remember(topicFor(m.Content)); ok {
		outStr = fmt.Sprintf("I heard %s", mem.orig)
		fired = true
	}

	return true, []string{outStr}
}

func Responder(botUser *discordgo.User) bot.Responder {
	recallMod := newRecallMod(botUser)
	return recallMod.Responder
}

func newRecallMod(botUser *discordgo.User) *recallMod {
	return &recallMod{}
}

func topicFor(content string) string {
	content = strings.TrimSuffix(content, "?")
	words := strings.Fields(content)

	minusNick := words[1:]
	return strings.Join(minusNick, " ")
}
