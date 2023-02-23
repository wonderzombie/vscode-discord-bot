package main

import (
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
