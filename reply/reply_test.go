package reply

import (
	"reflect"
	"testing"

	"github.com/wonderzombie/godiscbot/bot"
)

func Test_replyMod_Responder(t *testing.T) {
	helloWorld := []string{"hello world"}
	type fields struct {
		nick string
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
			fields{"foobot"},
			args{
				&bot.Message{Content: "hello foobot"},
			},
			true,
			[]string{"hello world"},
		},
		{
			"not mentioned",
			fields{"foobot"},
			args{
				&bot.Message{Content: "hello barbot"},
			},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &replyMod{
				nick:    tt.fields.nick,
				phrases: helloWorld,
			}
			gotFired, gotOut := rm.Responder(tt.args.m)
			if gotFired != tt.wantFired {
				t.Errorf("replyMod.Responder() got = %v, want %v", gotFired, tt.wantFired)
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("replyMod.Responder() got1 = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
