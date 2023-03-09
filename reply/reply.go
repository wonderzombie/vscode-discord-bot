package reply

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

func Responder(botUser *discordgo.User) bot.Responder {
	replyMod := newReplyMod(botUser)
	return replyMod.Responder
}

type replyMod struct {
	uniqueName string
	nick       string

	phrases []string
}

var defaultPhrases = []string{
	"who said that",
	"where is who, where is who",
}

func (rm *replyMod) Responder(m *bot.Message) (bool, []string) {
	lowered := strings.ToLower(m.Content)
	nick := strings.ToLower(rm.nick)
	if !strings.Contains(lowered, nick) {
		return false, nil
	}

	roll := rand.Int() % len(rm.phrases)
	out := rm.phrases[roll]

	return true, []string{out}
}

func newReplyMod(user *discordgo.User) *replyMod {
	return &replyMod{
		nick:       user.Username,
		uniqueName: user.String(),
		phrases:    defaultPhrases,
	}
}
