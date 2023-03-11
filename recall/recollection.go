package recall

import (
	"math/rand"
	"strings"
	"time"

	"github.com/wonderzombie/godiscbot/bot"
)

// a memory captures a message and when it was received.
type memory struct {
	orig string
	time.Time
}

// a recollection manages a list of memories
type recollection struct {
	memories []memory
}

// remember retrieves a message related to a topic (ostensibly a word or phrase) from the recollection's memories.
func (r *recollection) remember(topic string) (memory, bool) {
	var topicMemories []memory

	for _, m := range r.memories {
		if substr(m.orig, topic) {
			topicMemories = append(topicMemories, m)
		}
	}

	if len(topicMemories) == 0 {
		return memory{}, false
	}

	return choose(topicMemories), true
}

func (r *recollection) add(m *bot.Message, t time.Time) {
	r.memories = append(r.memories, memory{orig: m.Content, Time: t})
}

func choose(topicMemories []memory) memory {
	result := rand.Int() % len(topicMemories)
	return topicMemories[result]
}

func oldest(topicMemories []memory) memory {
	old := topicMemories[0]
	for _, m := range topicMemories {
		if older(m.Time, old.Time) {
			old = m
		}
	}
	return old
}

func older(older time.Time, newer time.Time) bool {
	return older.Local().Before(newer.Local())
}

func newest(topicMemories []memory) memory {
	var n memory
	for _, m := range topicMemories {
		if older(n.Time, m.Time) {
			n = m
		}
	}
	return n
}

func selectTopic(memories []memory, topic string) []memory {
	var topicMemories []memory

	for _, m := range memories {
		if substr(m.orig, topic) {
			topicMemories = append(topicMemories, m)
		}
	}

	return topicMemories
}

// substr returns true when substring is in phrase, and always ignores case.
func substr(phrase string, substring string) bool {
	lowerP := strings.ToLower(phrase)
	lowerSub := strings.ToLower(substring)
	return strings.Contains(lowerP, lowerSub)
}
