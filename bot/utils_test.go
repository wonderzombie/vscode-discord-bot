package bot

import "testing"

func TestEmpty(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{
			"empty slice",
			[]string{},
			true,
		},
		{
			"empty string in slice",
			[]string{""},
			true,
		},
		{
			"has strings",
			[]string{"hello", "world"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Empty(tt.lines); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
