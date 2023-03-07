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

type seenUser struct {
	username string
	t        time.Time
}

func New(users ...seenUser) bot.Responder {
	sm := &SeenModule{
		state: &seenState{
			seen: make(seenTimes, len(users)),
			mx:   sync.Mutex{},
		},
	}
	for _, u := range users {
		sm.state.seen[u.username] = u.t
	}
	return sm.Handle
}

func (sm *SeenModule) Handle(m *bot.Message) (bool, []string) {
	cmd, ok := m.Cmd()
	if !ok {
		return false, []string{}
	}

	var responder bot.Responder
	switch cmd {
	case "ping", "!ping":
		responder = sm.pong
	case "seen", "!seen":
		responder = sm.seenResp
	}

	return responder(m)
}

func (sm *SeenModule) pong(m *bot.Message) (bool, []string) {
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

func (sm *SeenModule) seenResp(m *bot.Message) (bool, []string) {
	lines := []string{}
	if !strings.HasPrefix(m.Content, "!seen") {
		return false, lines
	}

	sm.state.mx.Lock()
	defer sm.state.mx.Unlock()

	fields := strings.Fields(m.Content)
	var response string
	fired := false
	if len(fields) == 1 {
		response = sm.state.seen.String()
	} else if len(fields) == 2 {
		username := fields[1]
		if t, ok := sm.state.seen[username]; ok {
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
