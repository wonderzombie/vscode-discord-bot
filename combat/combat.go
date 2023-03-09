package combat

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/wonderzombie/godiscbot/bot"
)

const (
	defaultHp = 10
)

func randInt(n size) int {
	return rand.Int() % n.Int()
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

type quantity int
type size int

func (sz size) Int() int {
	return int(sz)
}

func (q quantity) Int() int {
	return int(q)
}

type tracker struct {
	m         map[string]*combatant
	defaultHp int
	roll      func(size) int
	rollN     func(size, quantity) int
}

func newTracker(defaultHp int, roller func(size) int) *tracker {
	return &tracker{
		m:         make(map[string]*combatant),
		defaultHp: defaultHp,
		roll:      roller,
		rollN: func(m size, q quantity) int {
			t := 0
			for i := 0; i < q.Int(); i++ {
				t += roller(m)
			}
			return t
		},
	}
}

// alwaysGet will get the combatant, creating a new entry for a combatant if needed.
func (cm *tracker) alwaysGet(k string) *combatant {
	if out, ok := cm.m[k]; ok {
		return out
	}
	return cm.init(k)
}

// init resets the combatant's status by creating a new record and returns the same.
func (cm *tracker) init(k string) *combatant {
	cm.m[k] = &combatant{
		name: k,
		hp:   defaultHp,
	}
	return cm.m[k]
}

func Responder() bot.Responder {
	tracker := newTracker(defaultHp, randInt)
	return tracker.Responder
}

func (cm *tracker) Responder(m *bot.Message) (fired bool, out []string) {
	msgParts := strings.Fields(m.Content)
	if len(msgParts) < 2 {
		return
	}

	action, targetName := msgParts[0], msgParts[1]
	target := cm.alwaysGet(targetName)
	out = cm.resolve(action, m.Author, target)
	if len(out) > 0 {
		fired = true
	}

	return fired, out
}

// https://go.dev/play/p/3R5wtH9yOMo
type resolver func(*tracker, string, *combatant) []string

func (cm *tracker) resolve(action string, author string, target *combatant) []string {
	var fn resolver = (*tracker).resolveNoop

	switch action := strings.Trim(action, "\n! "); action {
	case "heal", "heals", "bless", "cure", "aid":
		fn = (*tracker).resolveHeal
	case "res", "resurrect":
		fn = (*tracker).resolveRes
	case "attack", "hit", "stab", "bite", "curse":
		fn = (*tracker).resolveAttack
	}

	return fn(cm, author, target)
}

func (cm *tracker) resolveNoop(unused1 string, unused2 *combatant) []string {
	return nil
}

func (cm *tracker) resolveHeal(author string, target *combatant) []string {
	var outMsg string
	switch {
	case target.dead():
		outMsg = fmt.Sprintf("%s is dead!", target.name)
	case target.hp < defaultHp:
		healed := cm.rollN(2, 3)
		target.hp += healed
		outMsg = fmt.Sprintf("%s healed %s for %d hp!", author, target.name, healed)
	default:
		outMsg = fmt.Sprintf("%s already has %d hp!", target.name, target.hp)
	}

	return []string{outMsg}
}

func (cm *tracker) resolveRes(author string, target *combatant) []string {
	outMsg := fmt.Sprintf("%s is still alive!", target.name)
	if target.dead() {
		target.restore()
		outMsg = fmt.Sprintf("%s brings %s back from beyond the grave!", author, target.name)
	}
	return []string{outMsg}
}

func (cm *tracker) resolveAttack(author string, target *combatant) (out []string) {
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
