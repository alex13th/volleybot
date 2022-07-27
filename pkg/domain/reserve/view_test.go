package reserve

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

func TestTelegramView(t *testing.T) {
	plid, _ := uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	pl1 := person.Person{Id: plid, Firstname: "Elly"}
	person.NewPerson("Elly")
	tests := map[string]struct {
		res  Reserve
		text string
		str  string
	}{
		"2 hors": {
			res: Reserve{
				Person:    pl1,
				StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				Price:     600,
			},
			text: "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n💰 600 ₽",
			str:  "Сб, 04.12 15:00-17:00",
		},
		"With description": {
			res: Reserve{
				Person:      pl1,
				StartTime:   time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:     time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				Price:       600,
				Description: "Some description.",
			},
			text: "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n💰 600 ₽" +
				"\n\nSome description.",
			str: "Сб, 04.12 15:00-17:00",
		},
		"Canceled": {
			res: Reserve{
				Person:    pl1,
				StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				Canceled:  true,
			},
			text: "🔥 *ОТМЕНА* 🔥\n\n*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00",
			str:  "Сб, 04.12 15:00-17:00",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reserve := test.res
			tgv := NewTelegramViewRu(reserve)
			text := tgv.GetText()
			str := tgv.String()
			if text != test.text {
				t.Fail()
			}
			if str != test.str {
				t.Fail()
			}
		})
	}
}
