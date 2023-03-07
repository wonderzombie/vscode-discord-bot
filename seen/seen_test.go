package seen

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/wonderzombie/godiscbot/bot"
)

func messageCreate(username string, discriminator string, content string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:  user(username, discriminator),
			Content: content,
		}}
}

func message(username string, discriminator string, content string) *bot.Message {
	return bot.NewMessage(messageCreate(username, discriminator, content))
}

func user(username string, disc string) *discordgo.User {
	return &discordgo.User{
		Username:      username,
		Discriminator: "1111",
	}
}

func Test_pong(t *testing.T) {
	type args struct {
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
			// TODO: Add test cases.
			name: "testing ping",
			args: args{m: message("foo", "1111", "!ping")},

			ret: ret{
				wantStr:   []string{"PONG"},
				wantFired: true,
			},
		},
		{
			name: "testing pong",
			args: args{m: message("foo", "1111", "!pong")},
			ret: ret{
				wantStr:   []string{"PING"},
				wantFired: true,
			},
		},
		{
			name: "testing neither",
			args: args{m: message("foo", "1111", "!bees")},
			ret: ret{
				wantStr:   []string{""},
				wantFired: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := pong(tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("pong() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}

func Test_seen(t *testing.T) {
	initOnce.Do(initSeen)

	type args struct {
		m *bot.Message
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
			args: args{m: message("someuser", "1111", "!seen bar")},
			want: []string{"someuser#1111, I've never seen bar"},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"someuser#1111, I've never seen bar"},
			},
		},
		{
			name: "basic - !seen",
			args: args{m: message("foo", "1111", "!seen")},
			ret: ret{
				wantFired: true,
				wantStr:   []string{"NOBODY"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFired, gotStr := seenResp(tt.args.m); !reflect.DeepEqual(gotStr, tt.ret.wantStr) || gotFired != tt.ret.wantFired {
				t.Errorf("seen() = %v, %v; want %v, %v", gotFired, gotStr, tt.ret.wantFired, tt.ret.wantStr)
			}
		})
	}
}
