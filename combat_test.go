package main

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func fakeMap(name string, hp int) map[string]*combatant {
	return map[string]*combatant{
		name: {
			name: name,
			hp:   hp,
		},
	}
}

func Test_combatMap_get(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]*combatant
		k    string
		want *combatant
	}{
		{
			name: "present",
			m:    fakeMap("target", 11),
			k:    "target",
			want: &combatant{
				name: "target",
				hp:   11,
			},
		},
		{
			name: "empty",
			m:    map[string]*combatant{},
			k:    "target",
			want: &combatant{
				name: "target",
				hp:   defaultHp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m: tt.m,
			}
			if got := cm.get(tt.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combatMap.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combatMap_init(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]*combatant
		k    string
		want *combatant
	}{
		{
			name: "empty map",
			m:    map[string]*combatant{},
			k:    "target",
			want: &combatant{
				name: "target",
				hp:   defaultHp,
			},
		},
		{
			name: "existing entry initialized",
			m:    fakeMap("target", defaultHp*2),
			k:    "target",
			want: &combatant{
				name: "target",
				hp:   defaultHp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m: tt.m,
			}
			if got := cm.init(tt.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combatMap.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombat(t *testing.T) {
	type args struct {
		s *discordgo.Session
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Combat(tt.args.s, tt.args.m)
		})
	}
}

func Test_resolveAttack(t *testing.T) {
	type args struct {
		cmd    string
		author string
		target *combatant
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := resolveAttack(tt.args.cmd, tt.args.author, tt.args.target); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("resolveAttack() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_resolve(t *testing.T) {
	type args struct {
		cmd    string
		author string
		target *combatant
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := resolve(tt.args.cmd, tt.args.author, tt.args.target); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("resolve() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_resolveNoop(t *testing.T) {
	type args struct {
		unused1 string
		unused2 string
		unused3 *combatant
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEmpty := resolveNoop(tt.args.unused1, tt.args.unused2, tt.args.unused3); !reflect.DeepEqual(gotEmpty, tt.wantEmpty) {
				t.Errorf("resolveNoop() = %v, want %v", gotEmpty, tt.wantEmpty)
			}
		})
	}
}

func Test_resolveHeal(t *testing.T) {
	type args struct {
		cmd    string
		author string
		target *combatant
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveHeal(tt.args.cmd, tt.args.author, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveHeal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveRes(t *testing.T) {
	type args struct {
		cmd    string
		author string
		target *combatant
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveRes(tt.args.cmd, tt.args.author, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveRes() = %v, want %v", got, tt.want)
			}
		})
	}
}
