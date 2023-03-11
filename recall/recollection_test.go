package recall

import (
	"reflect"
	"testing"
	"time"
)

func Test_choose(t *testing.T) {
	type args struct {
		topicMemories []memory
	}
	tests := []struct {
		name string
		args args
		want memory
	}{
		{
			"one choice only",
			args{[]memory{{"foo", time.Unix(100, 0)}}},
			memory{"foo", time.Unix(100, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := choose(tt.args.topicMemories); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("choose() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oldest(t *testing.T) {
	type args struct {
		topicMemories []memory
	}
	tests := []struct {
		name string
		args args
		want memory
	}{
		{
			"two entries",
			args{
				[]memory{
					{"you are great", time.Unix(50000, 0)},
					{"you are not great", time.Unix(10000, 0)},
				},
			},
			memory{"you are not great", time.Unix(10000, 0)},
		},
		{
			"three entries",
			args{
				[]memory{
					{"you are something", time.Unix(90000, 0)},
					{"you are not great", time.Unix(50000, 0)},
					{"you are great", time.Unix(10000, 0)},
				},
			},
			memory{"you are great", time.Unix(10000, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oldest(tt.args.topicMemories); got.orig != tt.want.orig || !got.Time.Local().Equal(tt.want.Time.Local()) {
				t.Errorf("oldest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_older(t *testing.T) {
	type args struct {
		older time.Time
		newer time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"yes older",
			args{time.Unix(1000, 0).Local(), time.Unix(9000, 0)},
			true,
		},
		{
			"not older",
			args{time.Unix(9000, 0).Local(), time.Unix(1000, 0).Local()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := older(tt.args.older, tt.args.newer); got != tt.want {
				t.Errorf("older() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newest(t *testing.T) {
	type args struct {
		topicMemories []memory
	}
	tests := []struct {
		name string
		args args
		want memory
	}{
		{
			"two entries",
			args{
				[]memory{
					{"you are not great", time.Unix(50000, 0)},
					{"you are great", time.Unix(10000, 0)},
				},
			},
			memory{"you are not great", time.Unix(50000, 0)},
		},
		{
			"three entries",
			args{
				[]memory{
					{"you are something", time.Unix(90000, 0)},
					{"you are not great", time.Unix(50000, 0)},
					{"you are great", time.Unix(10000, 0)},
				},
			},
			memory{"you are something", time.Unix(90000, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newest(tt.args.topicMemories); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_substr(t *testing.T) {
	type args struct {
		phrase    string
		substring string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"lowercase, substr",
			args{"hello world", "world"},
			true,
		},
		{
			"mixed case, substr",
			args{"hElLo WoRlD", "world"},
			true,
		},
		{
			"lowercase, not substr",
			args{"hello world", "erroneous"},
			false,
		},
		{
			"mixed cased, not substr",
			args{"hElLo WoRlD", "erroneous"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := substr(tt.args.phrase, tt.args.substring); got != tt.want {
				t.Errorf("substr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectTopic(t *testing.T) {
	type args struct {
		memories []memory
		topic    string
	}
	tests := []struct {
		name string
		args args
		want []memory
	}{
		// TODO: Add test cases.
		{
			"only choice matches",
			args{
				[]memory{
					{"foobot, hello", time.Unix(1000, 0)},
				},
				"hello",
			},
			[]memory{
				{"foobot, hello", time.Unix(1000, 0)},
			},
		},
		{
			"one of three matches",
			args{
				[]memory{
					{"foobot, hello to the bees", time.Unix(1000, 0)},
					{"foobot, I love bees", time.Unix(2000, 0)},
					{"foobot, I love BATHS", time.Unix(2000, 0)},
				},
				"baths"},
			[]memory{
				{"foobot, I love BATHS", time.Unix(2000, 0)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectTopic(tt.args.memories, tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
