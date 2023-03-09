package reply

import (
	"fmt"
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

type replyMod struct {
	uniqueName string
	nick       string

	phrases []string
}

var defaultPhrases = []string{
	"who said that",
	"where is {{.Nick}}, where is {{.Nick}}",
}

func (rm *replyMod) Responder(m *bot.Message) (bool, []string) {
	lowered := strings.ToLower(m.Content)
	nick := strings.ToLower(rm.nick)
	if !strings.Contains(lowered, nick) {
		return false, nil
	}

	fmt.Println("[reply] saw my name mentioned")

	roll := rand.Int() % len(rm.phrases)
	out := rm.phrases[roll]

	return rm.renderReply(out)
}

func (rm *replyMod) renderReply(out string) (bool, []string) {
	tmpl, err := template.New("reply").Parse(out)
	if err != nil {
		fmt.Println("error loading template", tmpl.Name(), ":", err)
		return false, nil
	}

	buf := &strings.Builder{}
	err = tmpl.Execute(buf, struct{ Nick string }{rm.nick})
	if err != nil {
		fmt.Println("failed to execute template", tmpl.Name(), ":", err)
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
