package reserve

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

func TestTelegramView(t *testing.T) {
	plid, _ := uuid.Parse("da10db9a-490b-4010-9d8c-561cca979dd0")
	pl1 := person.Person{Id: plid, Firstname: "Elly"}
	person.NewPerson("Elly")
	plid, _ = uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	pl2 := person.Person{Id: plid, Firstname: "Steve"}
	plid, _ = uuid.Parse("80155587-168c-4255-82ec-991119f3e110")
	pl3 := person.Person{Id: plid, Firstname: "Tina", TelegramId: 123456}
	tests := map[string]struct {
		res  Reserve
		text string
		str  string
	}{
		"2 hors": {
			res: Reserve{
				Person:     pl1,
				StartTime:  time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				MinLevel:   int(Middle),
				MaxPlayers: 6,
				Price:      600,
			},
			text: "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n" +
				"💪*Уровень*: Средний\n💳 600\n*Игроков:* 6\n1.\n.\n.\n6.",
			str: "Сб, 04.12 15:00-17:00 (0/6)",
		},
		"4 max": {
			res: Reserve{
				Person:     pl1,
				StartTime:  time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				MaxPlayers: 4,
			},
			text: "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n" +
				"*Игроков:* 4\n1.\n2.\n3.\n4.",
			str: "Сб, 04.12 15:00-17:00 (0/4)",
		},
		"2 players": {
			res: Reserve{
				Person:     pl1,
				StartTime:  time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				MaxPlayers: 12,
				Players: map[uuid.UUID]Player{
					pl1.Id: {Person: pl1, Count: 2},
					pl2.Id: {Person: pl2, Count: 3},
					pl3.Id: {Person: pl3, Count: 1}}},
			text: "*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n" +
				"*Игроков:* 12\n1. Elly+1\n3. Steve+2\n6. [Tina](tg://user?id=123456)\n7.\n.\n.\n12.",
			str: "Сб, 04.12 15:00-17:00 (6/12)",
		},
		"Canceled": {
			res: Reserve{
				Person:     pl1,
				StartTime:  time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC),
				MaxPlayers: 12,
				Canceled:   true,
				Players: map[uuid.UUID]Player{
					pl1.Id: {Person: pl1, Count: 2},
				}},
			text: "🔥*ОТМЕНА*🔥\n*Elly*\n📆 Суббота, 04.12.2021\n⏰ 15:00-17:00\n" +
				"*Игроков:* 12\n1. Elly+1\n3.\n.\n.\n12.",
			str: "Сб, 04.12 15:00-17:00 (2/12)",
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
