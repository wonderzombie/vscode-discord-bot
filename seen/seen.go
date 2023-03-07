package seen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/wonderzombie/godiscbot/bot"
)

type seenTimes map[string]time.Time

func (st seenTimes) String() string {
	if len(st) == 0 {
		return "NOBODY"
	}

	var buf strings.Builder
	for user, t := range st {
		buf.WriteString(fmt.Sprintf("%v - %v\n", user, t.Format(time.UnixDate)))
	}
	return buf.String()
}

type seenState struct {
	seen seenTimes
	mx   sync.Mutex
}

type SeenModule struct {
	state *seenState
}

func NewSeenModule(ss *seenState) *bot.Module {
	sm := &SeenModule{ss}

	return bot.NewModule("seen",
		func(m *bot.Message) (bool, []string) {
			return sm.handleSeen(m)
		})
}

type SeenResponder func(ss *seenState, m *bot.Message) (bool, []string)

func (sm *SeenModule) handleSeen(m *bot.Message) (bool, []string) {
	cmd, ok := m.Cmd()
	if !ok {
		return false, []string{}
	}

	var responder SeenResponder
	switch cmd {
	case "ping", "!ping":
		responder = pong
	case "seen", "!seen":
		responder = seenResp
	}

	return responder(sm.state, m)
}

func pong(ss *seenState, m *bot.Message) (bool, []string) {
	var out string
	fired := false
	if m.Content == "!ping" {
		out = "PONG"
		fired = true
	} else if m.Content == "!pong" {
		out = "PING"
		fired = true
	}
	return fired, []string{out}
}

func seenResp(ss *seenState, m *bot.Message) (bool, []string) {
	lines := []string{}
	if !strings.HasPrefix(m.Content, "!seen") {
		return false, lines
	}

	ss.mx.Lock()
	defer ss.mx.Unlock()

	fields := strings.Fields(m.Content)
	var response string
	fired := false
	if len(fields) == 1 {
		response = ss.seen.String()
	} else if len(fields) == 2 {
		username := fields[1]
		if t, ok := ss.seen[username]; ok {
			response = fmt.Sprintf("%s, last time I saw %s it was %v", m.Author, username, t)
		} else {
			response = fmt.Sprintf("%s, I've never seen %s", m.Author, username)
		}
	}

	if response != "" {
		fired = true
		lines = append(lines, response)
	}
	return fired, lines
}
