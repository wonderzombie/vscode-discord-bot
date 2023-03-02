package bot

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func message(username string, discriminator string, content string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:  user(username, discriminator),
			Content: content,
		}}
}

func user(username string, disc string) *discordgo.User {
	return &discordgo.User{
		Username:      username,
		Discriminator: "1111",
	}
}

func Test_pong(t *testing.T) {
	type args struct {
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			// TODO: Add test cases.
			name: "testing ping",
			args: args{m: message("foo", "1111", "!ping")},
			want: []string{"PONG"},
		},
		{
			name: "testing pong",
			args: args{m: message("foo", "1111", "!pong")},
			want: []string{"PING"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pong(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seen(t *testing.T) {
	initOnce.Do(initSeen)

	type args struct {
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "basic - !seen foo",
			args: args{m: message("someuser", "1111", "!seen bar")},
			want: []string{"someuser#1111, I've never seen bar"},
		},
		{
			name: "basic - !seen",
			args: args{m: message("foo", "1111", "!seen")},
			want: []string{"NOBODY"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seenResp(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seen() = %v, want %v", got, tt.want)
			}
		})
	}
}
