package recall

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

type recallMod struct {
}

func (rm *recallMod) Responder(m *bot.Message) (bool, []string) {

	return false, nil
}

func Responder(botUser *discordgo.User) bot.Responder {
	recallMod := newRecallMod(botUser)
	return recallMod.Responder
}

func newRecallMod(botUser *discordgo.User) *recallMod {
	return &recallMod{}
}
