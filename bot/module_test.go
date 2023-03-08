package bot

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestMessage_Cmd(t *testing.T) {
	type args struct {
		Author  string
		Content string
		fields  []string
	}
	tests := []struct {
		name    string
		args    args
		wantCmd string
		wantOk  bool
	}{
		{
			"one arg command",
			args{"foo#1111", "!hello", []string{"!hello"}},
			"!hello",
			true,
		},
		{
			"zero arg command",
			args{"foo#1111", "hello!", []string{"hello!"}},
			"",
			false,
		},
		{
			"more than one arg command",
			args{"foo#1111", "!hello world :D", []string{"!hello", "world", ":D"}},
			"!hello",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Author:  tt.args.Author,
				Content: tt.args.Content,
				fields:  tt.args.fields,
			}
			got, got1 := m.Cmd()
			if got != tt.wantCmd {
				t.Errorf("Message.Cmd() got = %v, want %v", got, tt.wantCmd)
			}
			if got1 != tt.wantOk {
				t.Errorf("Message.Cmd() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func TestMessage_Args(t *testing.T) {
	type args struct {
		Author  string
		Content string
		fields  []string
	}
	tests := []struct {
		name      string
		fields    args
		wantOut   []string
		wantValid bool
	}{
		{
			"one arg",
			args{"foo#1111", "!seen bees", []string{"!seen", "bees"}},
			[]string{"bees"},
			true,
		},
		{
			"three args",
			args{"foo#1111", "!seen bees when", []string{"!seen", "bees", "when"}},
			[]string{"bees", "when"},
			true,
		},
		{
			"no args",
			args{"foo#1111", "!seen", []string{"!seen"}},
			[]string{""},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Author:  tt.fields.Author,
				Content: tt.fields.Content,
				fields:  tt.fields.fields,
			}
			gotOut, gotValid := m.Args()
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("Message.Args() gotOut = %q, want %q", gotOut, tt.wantOut)
			}
			if gotValid != tt.wantValid {
				t.Errorf("Message.Args() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func messageCreate(username string, discriminator string, content string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:  &discordgo.User{Username: username, Discriminator: discriminator},
			Content: content,
		}}
}

func TestNewMessage(t *testing.T) {
	type args struct {
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
		want *Message
	}{
		{
			"simple message",
			args{messageCreate("foo", "1111", "hello world")},
			&Message{Author: "foo#1111", Content: "hello world", fields: []string{"hello", "world"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessage(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
