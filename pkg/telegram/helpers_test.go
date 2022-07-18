package telegram

import (
	"testing"
	"time"
)

func TestBaseKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want BaseKeyboardHelper
	}{
		"Base data": {
			data: "state1_action1_somedata",
			want: BaseKeyboardHelper{State: "state1", Action: "action1", Data: "somedata"},
		},
		"Data with extras": {
			data: "state1_action1_some_data1",
			want: BaseKeyboardHelper{State: "state1", Action: "action1", Data: "some_data1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh, err := NewBaseKeyboardHelper(test.data, "_")
			test.want.sep = "_"

			if err != test.err {
				t.Fail()
			}

			if kh != test.want {
				t.Fail()
			}
		})
	}
}

func TestDateKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want DateKeyboardHelper
	}{
		"Set date action": {
			data: "date_set_2022-18-07_somedata",
			want: DateKeyboardHelper{
				BaseKeyboardHelper: BaseKeyboardHelper{State: "state1", Action: "set", Data: "somedata"},
				Days:               6, Columns: 2, DateFormat: "Mon, 02.01",
				Date: time.Date(2022, 07, 18, 0, 0, 0, 0, time.Local)},
		},
		"Notset date action": {
			data: "date_get_2022-18-07_somedata",
			want: DateKeyboardHelper{
				BaseKeyboardHelper: BaseKeyboardHelper{State: "state1", Action: "get", Data: "2022-18-07_somedata"},
				Days:               6, Columns: 2, DateFormat: "Mon, 02.01",
				Date: time.Time{}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh, err := NewDateKeyboardHelper(test.data, "_")
			test.want.sep = "_"

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
			data: "state1_action1_somedata",
			date: time.Date(2022, 07, 18, 0, 0, 0, 0, time.Local),
			want: "state1_set_2022-18-07_somedata",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh, err := NewDateKeyboardHelper(test.data, "_")

			if err != test.err {
				t.Fail()
			}
			if kh.GetBtnData(test.date) != test.want {
				t.Fail()
			}
		})
	}
}

func TestTimeKeyboardHelperParse(t *testing.T) {
	tests := map[string]struct {
		data string
		err  error
		want TimeKeyboardHelper
	}{
		"Set time action": {
			data: "date_set_22:10_somedata",
			want: TimeKeyboardHelper{
				BaseKeyboardHelper: BaseKeyboardHelper{State: "state1", Action: "set", Data: "somedata"},
				Time:               time.Date(0, 0, 0, 22, 10, 0, 0, time.Local)},
		},
		"Notset time action": {
			data: "date_get_22:10_somedata",
			want: TimeKeyboardHelper{
				BaseKeyboardHelper: BaseKeyboardHelper{State: "state1", Action: "get", Data: "22:10_somedata"},
				Time:               time.Time{}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh, err := NewTimeKeyboardHelper(test.data, "_")
			test.want.sep = "_"

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
			data: "state1_action1_somedata",
			time: time.Date(0, 0, 0, 22, 15, 0, 0, time.Local),
			want: "state1_set_22:15_somedata",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kh, err := NewTimeKeyboardHelper(test.data, "_")

			if err != test.err {
				t.Fail()
			}
			txt := kh.GetBtnData(test.time)
			if txt != test.want {
				t.Fail()
			}
		})
	}
}
