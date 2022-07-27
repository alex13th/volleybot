package telegram

import (
	"testing"
)

func TestStateParser(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want State
	}{
		"Minimal": {
			data: "pr_state1_action1",
			want: State{Prefix: "pr", State: "state1", Action: "action1"},
		},
		"With Data": {
			data: "pr_state1_action1_somedata",
			want: State{Prefix: "pr", State: "state1", Data: "somedata", Action: "action1"},
		},
		"With Value": {
			data: "pr_state1_action1_some_data1_dat",
			want: State{Prefix: "pr", State: "state1", Data: "some", Action: "action1", Value: "data1_dat"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewState().Parse(test.data)
			test.want.Separator = "_"
			if err != test.err {
				t.Fail()
			}

			if st != test.want {
				t.Fail()
			}
		})
	}
}

func TestStateString6(t *testing.T) {
	tests := map[string]struct {
		state State
		want  string
	}{
		"Minimal": {
			state: State{Prefix: "pr", State: "state1", Action: "action1"},
			want:  "pr_state1_action1",
		},
		"With Data": {
			state: State{Prefix: "pr", State: "state1", Data: "somedata", Action: "action1"},
			want:  "pr_state1_action1_somedata",
		},
		"With Value": {
			state: State{Prefix: "pr", State: "state1", Data: "some", Action: "action1", Value: "data1_dat"},
			want:  "pr_state1_action1_some_data1_dat",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.state.Separator = "_"
			text := test.state.String()
			if text != test.want {
				t.Fail()
			}
		})
	}
}
