package main

import (
	"fmt"
	"reflect"
	"testing"
)

func fakeMap(name string, hp int) map[string]*combatant {
	return map[string]*combatant{
		name: {
			name: name,
			hp:   hp,
		},
	}
}

func Test_combatMap_alwaysGet(t *testing.T) {
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
		{
			name: "dead",
			m:    fakeMap("target", -10),
			k:    "target",
			want: &combatant{
				name: "target",
				hp:   -10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m: tt.m,
			}
			if got := cm.alwaysGet(tt.k); !reflect.DeepEqual(got, tt.want) {
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

func Test_combatMap_resolveNoop(t *testing.T) {
	type fields struct {
		m map[string]*combatant
	}
	type args struct {
		unused1 string
		unused2 *combatant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"noop",
			fields{fakeMap("foo", 1)},
			args{"foo", &combatant{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m: tt.fields.m,
			}
			if gotEmpty := cm.resolveNoop(tt.args.unused1, tt.args.unused2); len(gotEmpty) != 0 {
				t.Errorf("combatMap.resolveNoop() = `%q`, want %q", gotEmpty, []string{})
			}
		})
	}
}

func rollMax(n int) int {
	return n
}

func nRollMax(q int, n int) int {
	return n * q
}

func Test_combatMap_resolveHeal(t *testing.T) {
	type fields struct {
		m     map[string]*combatant
		roll  func(int) int
		nRoll func(int, int) int
	}
	type args struct {
		author string
		target *combatant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			"heal damage",
			fields{fakeMap("hurted", defaultHp/2), rollMax, nRollMax},
			args{"authorperson", &combatant{"hurted", defaultHp / 2}},
			[]string{"authorperson healed hurted for 6 hp!"},
		},
		{
			"already full",
			fields{fakeMap("fullhealth", defaultHp), rollMax, nRollMax},
			args{"authorperson", &combatant{"fullhealth", defaultHp}},
			[]string{fmt.Sprintf("fullhealth already has %d hp!", defaultHp)},
		},
		{
			"dead",
			fields{fakeMap("deadguy", -defaultHp), rollMax, nRollMax},
			args{"authorperson", &combatant{"deadguy", -defaultHp}},
			[]string{"deadguy is dead!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m:     tt.fields.m,
				roll:  rollMax,
				nRoll: nRollMax,
			}
			if got := cm.resolveHeal(tt.args.author, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combatMap.resolveHeal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combatMap_resolveRes(t *testing.T) {
	type fields struct {
		m map[string]*combatant
	}
	type args struct {
		author string
		target *combatant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			"resurrect",
			fields{fakeMap("deadguy", -1)},
			args{"cleric", &combatant{"deadguy", -1}},
			[]string{"cleric brings deadguy back from beyond the grave!"},
		},
		{
			"failed resurrect",
			fields{fakeMap("aliveguy", 1)},
			args{"cleric", &combatant{"aliveguy", 1}},
			[]string{"aliveguy is still alive!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &combatMap{
				m:    tt.fields.m,
				roll: rollMax,
			}
			if got := cm.resolveRes(tt.args.author, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combatMap.resolveRes() = %v, want %v", got, tt.want)
			}
		})
	}
}
