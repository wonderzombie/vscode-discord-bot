package combat

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	defaultHp = 10
)

func randInt(n int) int {
	return rand.Int() % n
}

type combatant struct {
	name string
	hp   int
}

func (c *combatant) dead() bool {
	return c.hp < 1
}

func (c *combatant) restore() {
	c.hp = defaultHp
}

type combatMap struct {
	m         map[string]*combatant
	defaultHp int
	roll      func(int) int
	nRoll     func(int, int) int
}

func newCombatMap(defaultHp int, roller func(int) int) *combatMap {
	return &combatMap{
		m:         make(map[string]*combatant),
		defaultHp: defaultHp,
		roll:      roller,
		nRoll: func(q int, n int) int {
			t := 0
			for i := 0; i < q; i++ {
				t += roller(n)
			}
			return t
		},
	}
}

// alwaysGet will get the combatant, creating a new entry for a combatant if needed.
func (cm *combatMap) alwaysGet(k string) *combatant {
	if out, ok := cm.m[k]; ok {
		return out
	}
	return cm.init(k)
}

// init resets the combatant's status by creating a new record and returns the same.
func (cm *combatMap) init(k string) *combatant {
	cm.m[k] = &combatant{
		name: k,
		hp:   defaultHp,
	}
	return cm.m[k]
}

var (
	tracker *combatMap = newCombatMap(defaultHp, randInt)
	sent    []*discordgo.Message
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgParts := strings.Fields(m.Content)
	if len(msgParts) < 2 {
		return
	}

	action, targetName := msgParts[0], msgParts[1]
	target := tracker.alwaysGet(targetName)
	out := tracker.resolve(action, m.Author.Username, target)

	if len(out) == 0 {
		return
	}

	for _, msg := range out {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		cms, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Printf("failed to send message: %v\nMessage follows:\n%s", err, msg)
			continue
		}
		sent = append(sent, cms)
	}
}

// https://go.dev/play/p/3R5wtH9yOMo
type resolver func(*combatMap, string, *combatant) []string

func (cm *combatMap) resolve(action string, author string, target *combatant) []string {
	var fn resolver = (*combatMap).resolveNoop

	switch action := strings.Trim(action, "\n! "); action {
	case "heal", "heals", "bless", "cure", "aid":
		fn = (*combatMap).resolveHeal
	case "res", "resurrect":
		fn = (*combatMap).resolveRes
	case "attack", "hit", "stab", "bite", "curse":
		fn = (*combatMap).resolveAttack
	}

	return fn(cm, author, target)
}

func (cm *combatMap) resolveNoop(unused1 string, unused2 *combatant) []string {
	return []string{}
}

func (cm *combatMap) resolveHeal(author string, target *combatant) []string {
	var outMsg string
	switch {
	case target.dead():
		outMsg = fmt.Sprintf("%s is dead!", target.name)
	case target.hp < defaultHp:
		healed := cm.nRoll(2, 3)
		target.hp += healed
		outMsg = fmt.Sprintf("%s healed %s for %d hp!", author, target.name, healed)
	default:
		outMsg = fmt.Sprintf("%s already has %d hp!", target.name, target.hp)
	}

	return []string{outMsg}
}

func (cm *combatMap) resolveRes(author string, target *combatant) []string {
	outMsg := fmt.Sprintf("%s is still alive!", target.name)
	if target.dead() {
		target.restore()
		outMsg = fmt.Sprintf("%s brings %s back from beyond the grave!", author, target.name)
	}
	return []string{outMsg}
}

func (cm *combatMap) resolveAttack(author string, target *combatant) (out []string) {
	// can't attack dead people
	if target.dead() {
		out = []string{fmt.Sprintf("%s is already dead!", target.name)}
		return out
	}

	// resolve attack: take damage 9 in 10 times
	out = []string{fmt.Sprintf("%s's attack misses %s!", author, target.name)}
	if attackRoll := cm.roll(10); attackRoll > 1 {
		damage := cm.roll(6)
		target.hp -= damage
		out = []string{fmt.Sprintf("%s hits %s for %d damage!", author, target.name, damage)}

		// if target has no hp remaining
		if target.dead() {
			out = append(out, fmt.Sprintf("%s dies!", target.name))
		}
	}

	return out
}
