package views

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
)

func TestTelegramView(t *testing.T) {
	tests := map[string]struct {
		start time.Time
		end   time.Time
		want  string
	}{
		"2 hors": {
			start: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
			want:  "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, _ := person.NewPerson("Elly")
			reserve, _ := reserve.NewReserve(p, test.start, test.end)
			tgv := NewTelegramViewRu(reserve)
			text := tgv.GetText()
			if text != test.want {
				t.Fail()
			}
		})
	}
}
