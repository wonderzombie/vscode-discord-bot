package reply

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

func Responder(botUser *discordgo.User) bot.Responder {
	replyMod := newReplyMod(botUser)
	return replyMod.Responder
}

func Responder2(botUser *discordgo.User) bot.Responder2 {
	replyMod := newReplyMod(botUser)
	return replyMod.Responder2
}

type replyMod struct {
	nick, uniqueName string
	phrases          []string
}

func (rm *replyMod) Responder2(w io.StringWriter, m *bot.Message) (fired bool) {
	if fired, out := rm.Responder(m); fired {
		for _, o := range out {
			w.WriteString(o)
		}
	}
	return fired
}

func (rm *replyMod) Responder(m *bot.Message) (bool, []string) {
	if m == nil {
		return false, nil
	}

	lowered := strings.ToLower(m.Content)
	nick := strings.ToLower(rm.nick)
	if !strings.Contains(lowered, nick) {
		return false, nil
	}

	fmt.Println("[reply] saw my name mentioned")

	roll := rand.Int() % len(rm.phrases)
	out := rm.phrases[roll]

	return rm.renderReply(out, m)
}

type tmplArgs struct {
	Phrase string // the phrase (reply) which has been chosen
	Nick   string // the nickname of the bot user
	Author string // author of the message sans discriminator
}

var defaultPhrases = []string{
	// TODO: consider an exercise to add real mentions using a sub-template.
	"who said that? was it you, {{.Author}}?",
	"where is {{.Nick}}, where is {{.Nick}}",
}

func (rm *replyMod) renderReply(rawPhrase string, m *bot.Message) (bool, []string) {
	if m == nil || rawPhrase == "" {
		return false, nil
	}

	tmpl, err := template.New("reply").Parse(rawPhrase)
	if err != nil {
		fmt.Println("error loading template:", err)
		return false, nil
	}

	buf := &strings.Builder{}
	err = tmpl.Execute(buf, tmplArgs{rawPhrase, rm.nick, m.Author})
	if err != nil {
		fmt.Println("failed to execute template:", err)
		return false, nil
	}
	return true, []string{buf.String()}
}

func newReplyMod(user *discordgo.User) *replyMod {
	return &replyMod{
		nick:       user.Username,
		uniqueName: user.String(),
		phrases:    defaultPhrases,
	}
}
