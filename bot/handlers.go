package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	handlers = []MessageHandler{
		seeing,
	}

	readiers = []ReadyHandler{}

	responders = []MessageResponder{
		pong,
		sent,
	}
)

type seenTimes map[string]time.Time

func (st seenTimes) String() string {
	if len(st) == 0 {
		return ""
	}

	var buf strings.Builder
	for user, t := range st {
		buf.WriteString(fmt.Sprintf("%v - %v\n", user, t.Format(time.UnixDate)))
	}
	return buf.String()
}

type SeenState struct {
	user *discordgo.User
	seen seenTimes
	mx   sync.Mutex
	sent []*discordgo.Message
}

func NewSeen() *SeenState {
	return &SeenState{
		seen: make(seenTimes),
		mx:   sync.Mutex{},
		sent: []*discordgo.Message{},
		// user: s.State.User,
	}
}

func (ss *SeenState) addSent(out *discordgo.Message) {
	ss.sent = append(ss.sent, out)
}

var currentState *SeenState

func initSeen() {
	currentState = NewSeen()
}

// initOnce is an experiment. init() is called before main, but in tests, there's no main per se. To ensure that
// there's no double-calls to initSeen, both this and tests use initOnce.
var initOnce sync.Once

func init() {
	initOnce.Do(initSeen)
}

func Seen(s *discordgo.Session, m *discordgo.MessageCreate) {
	var responder MessageResponder

	fields := strings.Fields(m.Content)
	if empty(fields) {
		return
	}

	switch cmd := fields[0]; cmd {
	case "!ping":
		responder = pong
	case "!seen":
		responder = seen
	}

	responses := responder(m)
	if len(responses) == 0 {
		return
	}

	for _, resp := range responses {
		m, err := s.ChannelMessageSend(m.ChannelID, resp)
		if err != nil {
			log.Printf("WARNING: unable to send: %v\nMessage content: %v\n", err, resp)
		}
		currentState.addSent(m)
	}
}

func pong(m *discordgo.MessageCreate) []string {
	var out string
	if m.Content == "!ping" {
		out = "PONG"
	} else if m.Content == "!pong" {
		out = "PING"
	}
	return []string{out}
}

func seen(m *discordgo.MessageCreate) []string {
	lines := []string{}
	if !strings.HasPrefix(m.Content, "!seen") {
		return lines
	}

	currentState.mx.Lock()
	defer currentState.mx.Unlock()

	fields := strings.Fields(m.Content)
	var response string
	if len(fields) == 1 {
		response = fmt.Sprintf("%q", currentState.seen)
	} else if len(fields) == 2 {
		username := fields[1]
		if t, ok := currentState.seen[username]; ok {
			response = fmt.Sprintf("%s, I saw %s and it was %v", m.Author, username, t)
		} else {
			response = fmt.Sprintf("%s, I've never seen %s", m.Author, username)
		}
	}

	if response != "" {
		lines = append(lines, response)
	}
	return lines
}

// TODO: something like a list. Maybe confined to an allow list of channels in .env or similar.
func seeing(m *discordgo.MessageCreate) {
	currentState.mx.Lock()
	defer currentState.mx.Unlock()
	if _, ok := currentState.seen[m.Author.Username]; !ok {
		currentState.seen[m.Author.Username] = m.Timestamp
	}
}

func sent(m *discordgo.MessageCreate) []string {
	lines := []string{}
	if !strings.HasPrefix(m.Content, "!sent") || len(currentState.sent) == 0 {
		return lines
	}

	var buf strings.Builder
	buf.WriteString("I sent these:\n")
	for i, msg := range currentState.sent {
		out := fmt.Sprintf("%d %s %s %s\n", i, msg.ChannelID, msg.Timestamp, msg.Content)
		buf.WriteString(out)
	}
	lines = append(lines, buf.String())
	return lines
}
