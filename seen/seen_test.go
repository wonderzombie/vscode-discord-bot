package seen

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

func message(username string, discriminator string, content string) *bot.Message {
	return bot.NewMessage(&discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:  user(username, discriminator),
			Content: content,
		}})
}

func user(username string, disc string) *discordgo.User {
	return &discordgo.User{
		Username:      username,
		Discriminator: "1111",
	}
}

func Test_pong(t *testing.T) {
	emptySeen := initSeen()
	seenMod := SeenModule{emptySeen}

	type args struct {
		r bot.Responder
		m *bot.Message
	}
	type ret struct {
		wantStr   []string
		wantFired bool
	}
	tests := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "testing ping",
			args: args{
				m: message("foo", "1111", "!ping"),
				r: seenMod.pong,
			},
			ret: ret{
				wantStr:   []string{"PONG"},
				wantFired: true,
			},
		},
		{
			name: "testing pong",
			args: args{
				m: message("foo", "1111", "!pong"),
				r: seenMod.pong,
			},
			ret: ret{
				wantStr:   []string{"PING"},
				wantFired: true,
			},
		},
		{
			name: "testing neither",
			args: args{
				m: message("foo", "1111", "!bees"),
				r: seenMod.pong,
			},
			ret: ret{
				wantStr:   []string{""},
				wantFired: false,
			},
		},
	}
	for _, tt := range tests {
		responder := tt.args.r
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := responder(tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("pong() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}

func initSeen(users ...seenUser) *seenState {
	out := make(seenTimes, len(users))
	for _, u := range users {
		out[u.username] = u.t
	}
	return &seenState{
		seen: out,
		mx:   sync.Mutex{},
	}

}

func Test_seen(t *testing.T) {
	type args struct {
		r bot.Responder
		m *bot.Message
	}
	type ret struct {
		wantFired bool
		wantStr   []string
	}
	tests := []struct {
		name string
		args args
		ret  ret
	}{
		{
			name: "never seen bar",
			args: args{
				r: New(seenUser{"not_bar", time.Unix(1, 0)}),
				m: message("someuser", "1111", "!seen bar")},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"someuser#1111, I've never seen bar"},
			},
		},
		{
			name: "nobody seen",
			args: args{
				r: New(),
				m: message("foo", "1111", "!seen")},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"NOBODY"},
			},
		},
		{
			name: "someone seen",
			args: args{
				m: message("foo", "1111", "!seen bar"),
				r: New(seenUser{"bar", time.Unix(1, 0)}),
			},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"foo#1111, last time I saw bar it was 1969-12-31 16:00:01 -0800 PST"},
			},
		},
	}
	for _, tt := range tests {
		sm := tt.args.r
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := sm(tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("seen() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}
