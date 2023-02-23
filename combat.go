package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	defaultHp = 10
)

type combatant struct {
	name string
	hp   int
}

type combatMap struct {
	m map[string]*combatant
}

func (cm *combatMap) get(k string) *combatant {
	if _, ok := cm.m[k]; !ok {
		cm.m[k] = &combatant{
			name: k,
			hp:   defaultHp,
		}
	}
	return cm.m[k]
}

var tracker = combatMap{}

var combatVerbs = map[string]bool{
	"!attack": true,
	// "!heal":   true,
	// "!res":    true,
}

func Combat(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgParts := strings.Fields(m.Content)
	if len(msgParts) < 2 {
		return
	}

	cmd, targetName := msgParts[0], msgParts[1]
	if !combatVerbs[cmd] {
		return
	}

	target := tracker.get(targetName)
	out := resolve(cmd, m.Author.Username, target)

	if len(out) == 0 {
		return
	}

	for _, msg := range out {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func resolve(cmd string, author string, target *combatant) (out []string) {
	// TODO: change to a switch stmt
	if !strings.HasSuffix(cmd, "attack") {
		return
	}

	if target.hp <= 0 {
		out = []string{fmt.Sprintf("%s is already dead!", target.name)}
		return out
	}

	// resolving an attack
	if attackRoll := rand.Int() % 10; attackRoll > 1 {
		damage := rand.Int() % 6
		target.hp -= damage
		out = []string{fmt.Sprintf("%s hits %s for %d damage!", author, target.name, damage)}
	} else {
		out = []string{fmt.Sprintf("%s's attack misses %s!", author, target.name)}
	}

	if target.hp >= 0 {
		return out
	}

	// resolving death
	out = append(out, fmt.Sprintf("%s dies!", target.name))

	return out
}
