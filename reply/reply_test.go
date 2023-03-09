package reply

import (
	"reflect"
	"testing"

	"github.com/wonderzombie/godiscbot/bot"
)

func Test_replyMod_Responder(t *testing.T) {
	helloWorld := []string{"hello world"}
	type fields struct {
		nick    string
		phrases []string
	}
	type args struct {
		m *bot.Message
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFired bool
		wantOut   []string
	}{
		{
			"mentioned",
			fields{"foobot", helloWorld},
			args{
				&bot.Message{Content: "hello foobot"},
			},
			true,
			[]string{"hello world"},
		},
		{
			"not mentioned",
			fields{"foobot", helloWorld},
			args{
				&bot.Message{Content: "hello barbot"},
			},
			false,
			nil,
		},
		{
			"reply has nick",
			fields{"foobot", []string{"i am {{.Nick}}!"}},
			args{
				&bot.Message{Content: "what's up foobot"},
			},
			true,
			[]string{"i am foobot!"},
		},
		{
			"reply has user's nick",
			fields{"foobot", []string{"hello {{.Author}}"}},
			args{
				&bot.Message{Content: "what's up foobot", Author: "someuser"},
			},
			true,
			[]string{"hello someuser"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &replyMod{
				nick:    tt.fields.nick,
				phrases: tt.fields.phrases,
			}
			gotFired, gotOut := rm.Responder(tt.args.m)
			if gotFired != tt.wantFired {
				t.Errorf("replyMod.Responder() gotFired = %v, want %v", gotFired, tt.wantFired)
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("replyMod.Responder() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
