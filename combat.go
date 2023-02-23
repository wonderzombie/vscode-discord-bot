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

// https://go.dev/play/p/3R5wtH9yOMo
type resolver func(*combatMap, string, *combatant) []string

var tracker = combatMap{}

func Combat(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgParts := strings.Fields(m.Content)
	if len(msgParts) < 2 {
		return
	}

	cmd, targetName := msgParts[0], msgParts[1]
	target := tracker.get(targetName)
	out := tracker.resolve(cmd, m.Author.Username, target)

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

func (cm *combatMap) resolve(cmd string, author string, target *combatant) (out []string) {
	// TODO: change to a switch stmt
	var fn resolver = (*combatMap).resolveNoop

	switch cmd := strings.Trim(cmd, "\n! "); cmd {
	case "heal", "heals", "bless", "cure", "aid":
		fn = (*combatMap).resolveHeal
	case "res", "resurrect":
		fn = (*combatMap).resolveRes
	case "attack", "hit", "stab", "bite", "curse":
		fn = (*combatMap).resolveAttack
	}

	return fn(cm, author, target)
}

func (cm *combatMap) resolveNoop(unused1 string, unused2 *combatant) (empty []string) {
	return empty
}

func (cm *combatMap) resolveHeal(author string, target *combatant) []string {
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

func (cm *combatMap) resolveRes(author string, target *combatant) []string {
	var outMsg string
	if target.hp > 0 {
		outMsg = fmt.Sprintf("%s is still alive, with %d hp.", target.name, target.hp)
		return []string{outMsg}
	}

	target.hp = defaultHp
	outMsg = fmt.Sprintf("%s brings %s back from beyond the grave!", author, target.name)
	return []string{outMsg}
}

func (cm *combatMap) resolveAttack(author string, target *combatant) (out []string) {
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
