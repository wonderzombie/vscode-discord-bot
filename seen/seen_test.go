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

	type args struct {
		m  *bot.Message
		ss *seenState
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
				m:  message("foo", "1111", "!ping"),
				ss: emptySeen,
			},
			ret: ret{
				wantStr:   []string{"PONG"},
				wantFired: true,
			},
		},
		{
			name: "testing pong",
			args: args{
				m:  message("foo", "1111", "!pong"),
				ss: emptySeen,
			},
			ret: ret{
				wantStr:   []string{"PING"},
				wantFired: true,
			},
		},
		{
			name: "testing neither",
			args: args{
				m:  message("foo", "1111", "!bees"),
				ss: emptySeen,
			},
			ret: ret{
				wantStr:   []string{""},
				wantFired: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := pong(tt.args.ss, tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("pong() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}

type seenUser struct {
	name string
	t    time.Time
}

func initSeen(users ...seenUser) *seenState {
	out := make(seenTimes, len(users))
	for _, su := range users {
		out[su.name] = su.t
	}
	return &seenState{
		seen: out,
		mx:   sync.Mutex{},
	}

}

func Test_seen(t *testing.T) {
	type args struct {
		ss *seenState
		m  *bot.Message
	}
	type ret struct {
		wantFired bool
		wantStr   []string
	}
	tests := []struct {
		name string
		args args
		want []string
		ret  ret
	}{
		{
			name: "basic - !seen foo",
			args: args{
				m:  message("someuser", "1111", "!seen bar"),
				ss: initSeen(seenUser{name: "baz", t: time.Unix(1, 0)})},
			want: []string{"someuser#1111, I've never seen bar"},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"someuser#1111, I've never seen bar"},
			},
		},
		{
			name: "basic - !seen",
			args: args{
				m:  message("foo", "1111", "!seen"),
				ss: initSeen()},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"NOBODY"},
			},
		},
		{
			name: "basic - !seen someone",
			args: args{
				m:  message("foo", "1111", "!seen bar"),
				ss: initSeen(seenUser{name: "bar", t: time.Unix(1, 0)}),
			},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"foo#1111, last time I saw bar it was 1969-12-31 16:00:01 -0800 PST"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := seenResp(tt.args.ss, tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("seen() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}
