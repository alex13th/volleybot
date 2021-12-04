package reserve

import (
	"testing"
	"time"
)

func TestReserveGetDuration(t *testing.T) {
	tests := map[string]struct {
		start time.Time
		end   time.Time
		want  float64
	}{
		"2 hors": {
			start: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
			want:  120,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reserve := Reserve{StartTime: test.start, EndTime: test.end}
			duration := reserve.GetDuration().Minutes()
			if duration != test.want {
				t.Fail()
			}
		})
	}
}

func TestReserveCheckConflicts(t *testing.T) {
	reserve := Reserve{
		StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
	}

	tests := map[string]struct {
		other Reserve
		want  bool
	}{
		"Other before": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 13, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		"Other after": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 19, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		"Other in": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 16, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		"Other over": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 19, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		"Other before in": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 16, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		"Other in after": {
			other: Reserve{
				StartTime: time.Date(2021, 12, 04, 16, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 20, 0, 0, 0, time.UTC),
			},
			want: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if reserve.CheckConflicts(test.other) != test.want {
				t.Fail()
			}
		})
	}
}
