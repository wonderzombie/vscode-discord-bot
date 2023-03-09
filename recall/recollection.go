package recall

import (
	"math/rand"
	"strings"
	"time"
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

type memoryOpt int

const (
	RANDOM memoryOpt = iota
	OLDEST
	NEWEST
)

func (r *recollection) remember(topic string, opt memoryOpt) memory {
	var topicMemories []memory

	for _, m := range r.memories {
		if substr(m.orig, topic) {
			topicMemories = append(topicMemories, m)
		}
	}

	var fn func([]memory) memory
	switch opt {
	case RANDOM:
		fn = choose
	case OLDEST:
		fn = oldest
	case NEWEST:
		fn = newest
	}

	return fn(topicMemories)
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
