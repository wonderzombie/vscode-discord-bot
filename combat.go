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
		cm.init(k)
	}
	return cm.m[k]
}

func (cm *combatMap) init(k string) *combatant {
	cm.m[k] = &combatant{
		name: k,
		hp:   defaultHp,
	}
	return cm.m[k]
}

var tracker = combatMap{}

var combatVerbs = map[string]bool{
	"!attack": true,
	// "!heal":   true,
	// "!res":    true,
}

type resolver func(string, string, *combatant) []string

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

func resolveAttack(cmd string, author string, target *combatant) (out []string) {
	// TODO: change to a switch stmt
	if !strings.HasSuffix(cmd, "attack") {
		return
	}

	// can't attack dead people
	if target.hp <= 0 {
		out = []string{fmt.Sprintf("%s is already dead!", target.name)}
		return out
	}

	// resolve attack: one is a miss, else take damage
	if attackRoll := rand.Int() % 10; attackRoll > 1 {
		damage := rand.Int() % 6
		target.hp -= damage
		out = []string{fmt.Sprintf("%s hits %s for %d damage!", author, target.name, damage)}
	} else {
		out = []string{fmt.Sprintf("%s's attack misses %s!", author, target.name)}
	}

	// target has hp remaining
	if target.hp >= 0 {
		return out
	}

	// target has no hp remaining
	out = append(out, fmt.Sprintf("%s dies!", target.name))

	return out

}

func resolve(cmd string, author string, target *combatant) (out []string) {
	// TODO: change to a switch stmt
	var fn resolver = resolveNoop
	if strings.HasSuffix(cmd, "attack") {
		fn = resolveAttack
	}

	switch cmd {
	case "!heal":
		fn = resolveHeal
	case "!res", "!resurrect":
		fn = resolveRes
	case "!attack", "!hit":
		fn = resolveAttack
	}

	return fn(cmd, author, target)
}

func resolveNoop(unused1 string, unused2 string, unused3 *combatant) (empty []string) {
	return empty
}

func resolveHeal(cmd string, author string, target *combatant) []string {
	// TODO: change to a switch stmt
	if !strings.HasSuffix(cmd, "heal") {
		return []string{}
	}
	outMsg := ""

	if target.hp > defaultHp {
		outMsg = fmt.Sprintf("%s already has %d hp!", target.name, target.hp)
	} else {
		healed := rand.Int() % 8
		target.hp += healed
		outMsg = fmt.Sprintf("%s healed %s for %d hp!", author, target.name, target.hp)
	}
	return []string{outMsg}
}

func resolveRes(cmd string, author string, target *combatant) []string {
	if !strings.HasSuffix(cmd, "res") || !strings.HasSuffix(cmd, "resurrect") {
		return []string{}
	}

	var outMsg string
	if target.hp > 0 {
		outMsg = fmt.Sprintf("%s is still alive, with %d hp.", target.name, target.hp)
		return []string{outMsg}
	}

	target.hp = defaultHp
	outMsg = fmt.Sprintf("%s brings %s back from beyond the grave!", author, target.name)
	return []string{outMsg}
}
