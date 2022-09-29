package telegram

import (
	"strconv"
	"testing"
	"time"
)

func TestBaseCallbackDataParser(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want State
	}{
		"Base data": {
			data: "pr_state1_action1_somedata",
			want: State{Prefix: "pr", State: "state1", Data: "somedata", Action: "action1"},
		},
		"Data with extras": {
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

func TestDateKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want struct {
			Data string
			Date time.Time
		}
	}{
		"Set date action": {
			data: "pr_date_set_somedata_2022-07-18",
			want: struct {
				Data string
				Date time.Time
			}{
				Data: "somedata",
				Date: time.Date(2022, 07, 18, 0, 0, 0, 0, time.Local),
			},
		},
		"Notset date action": {
			data: "pr_date_get_somedata_2022-07-18",
			want: struct {
				Data string
				Date time.Time
			}{
				Data: "somedata",
				Date: time.Time{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh := NewDateKeyboardHelper()
			st, _ := NewState().Parse(test.data)
			kh.State = st
			err := kh.Parse()

			if err != test.err {
				t.Fail()
			}

			if !kh.Date.Equal(test.want.Date) {
				t.Fail()
			}
			if kh.Data != test.want.Data {
				t.Fail()
			}
		})
	}
}

func TestDateKeyboardHelperGetBtnData(t *testing.T) {
	tests := map[string]struct {
		data string
		date time.Time
		err  error
		want string
	}{
		"2022-18-07": {
			data: "pr_state1_action1_somedata",
			date: time.Date(2022, 07, 18, 0, 0, 0, 0, time.Local),
			want: "pr_state1_set_somedata_2022-07-18",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewState().Parse(test.data)

			if err != test.err {
				t.Fail()
			}
			st.Value = test.date.Format("2006-01-02")
			st.Action = "set"
			act := st.String()
			if act != test.want {
				t.Fail()
			}
		})
	}
}

func TestTimeKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want struct {
			Data string
			Time time.Time
		}
	}{
		"Set time action": {
			data: "pr_date_set_somedata_22:10",
			want: struct {
				Data string
				Time time.Time
			}{
				Data: "somedata",
				Time: time.Date(0, 0, 0, 22, 10, 0, 0, time.Local),
			},
		},
		"Notset time action": {
			data: "pr_date_get_somedata_22:10",
			want: struct {
				Data string
				Time time.Time
			}{
				Data: "somedata",
				Time: time.Time{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh := NewTimeKeyboardHelper()
			st, _ := NewState().Parse(test.data)
			kh.State = st
			err := kh.Parse()

			if err != test.err {
				t.Fail()
			}
			if !kh.Time.Equal(test.want.Time) {
				t.Fail()
			}
			if kh.Data != test.want.Data {
				t.Fail()
			}
		})
	}
}

func TestTimeKeyboardHelperGetBtnData(t *testing.T) {
	tests := map[string]struct {
		data string
		time time.Time
		err  error
		want string
	}{
		"22:15": {
			data: "pr_state1_action1_somedata",
			time: time.Date(0, 0, 0, 22, 15, 0, 0, time.Local),
			want: "pr_state1_set_somedata_22:15",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewState().Parse(test.data)

			if err != test.err {
				t.Fail()
			}
			st.Value = test.time.Format("15:04")
			st.Action = "set"
			txt := st.String()
			if txt != test.want {
				t.Fail()
			}
		})
	}
}

func TestActionsKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want struct {
			State  string
			Action string
		}
	}{
		"Action1": {
			data: "pr_state1_action1",
			want: struct {
				State  string
				Action string
			}{State: "state1", Action: "action1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh := NewActionsKeyboardHelper()
			st, _ := NewState().Parse(test.data)
			kh.State = st

			if !(kh.Action == test.want.Action) {
				t.Fail()
			}
		})
	}
}

func TestActionsKeyboardHelperGetBtnData(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want string
	}{
		"action1": {
			data: "pr_state1_action0_somedata",
			want: "pr_state1_action1_somedata",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewState().Parse(test.data)
			if err != nil {
				t.Fail()
			}
			st.Action = "action1"
			txt := st.String()
			if txt != test.want {
				t.Fail()
			}
		})
	}
}

func TestCountKeyboardHelperGetBtnData(t *testing.T) {
	tests := map[string]struct {
		data  string
		count int
		err   error
		want  string
	}{
		"12": {
			data:  "pr_state1_action1_somedata",
			count: 12,
			want:  "pr_state1_set_somedata_12",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewState().Parse(test.data)
			if err != test.err {
				t.Fail()
			}
			st.Value = strconv.Itoa(test.count)
			st.Action = "set"
			txt := st.String()
			if txt != test.want {
				t.Fail()
			}
		})
	}
}
