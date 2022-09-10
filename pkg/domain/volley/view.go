package volley

import (
	"fmt"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/goodsign/monday"
)

type TelegramViewResources struct {
	reserve.TelegramViewResources
	GameLabel       string
	TrainingLabel   string
	TournamentLabel string
	TennisLabel     string
}

func NewTelegramResourcesRu() TelegramViewResources {
	rres := reserve.NewTelegramResourcesRu()
	return TelegramViewResources{TelegramViewResources: rres,
		GameLabel:       "🏐 *СВОБОДНЫЕ ИГРЫ* 🏐",
		TrainingLabel:   "‼️ *ТРЕНИРОВКА* ‼️",
		TournamentLabel: "💥🔥 *ТУРНИР* 🔥💥",
		TennisLabel:     "🎾 *ПЛЯЖНЫЙ ТЕННИС* 🎾",
	}
}

func NewTelegramViewRu(v Volley) TelegramView {
	return TelegramView{Volley: v, TelegramViewResources: NewTelegramResourcesRu()}
}

type TelegramView struct {
	Volley
	TelegramViewResources
}

func (tgv *TelegramView) String() string {
	plcount := 0
	for _, pl := range tgv.Volley.Players {
		plcount += pl.Count
	}
	return fmt.Sprintf("%s %s %s (%d/%d)", tgv.Volley.Activity.Emoji(),
		monday.Format(tgv.Reserve.StartTime, "Mon, 02.01", tgv.Locale),
		tgv.GetTimeText(), plcount, tgv.Volley.MaxPlayers)
}

func (tgv *TelegramView) GetText() (text string) {
	if tgv.Volley.Canceled {
		text = tgv.CancelLabel + "\n\n"
	} else if tgv.Volley.Activity == 10 {
		text = tgv.TrainingLabel + "\n\n"
	} else if tgv.Volley.Activity == 20 {
		text = tgv.TournamentLabel + "\n\n"
	} else if tgv.Volley.Activity == 30 {
		text = tgv.TennisLabel + "\n\n"
	} else {
		text = tgv.GameLabel + "\n\n"
	}
	var uname string
	if tgv.Reserve.Person.TelegramId != 0 {
		uname = fmt.Sprintf("[%s](tg://user?id=%d)", tgv.Reserve.Person.String(), tgv.Reserve.Person.TelegramId)
	} else {
		uname = fmt.Sprintf("*%s*", tgv.Reserve.Person.String())
	}
	text += fmt.Sprintf("%s\n%s %s\n%s %s", uname,
		tgv.DateLabel, monday.Format(tgv.Reserve.StartTime, "Monday, 02.01.2006", tgv.Locale),
		tgv.TimeLabel, tgv.GetTimeText())
	if tgv.Volley.MinLevel > 0 {
		text += fmt.Sprintf("\n💪*Уровень*: %s", person.PlayerLevel(tgv.Volley.MinLevel))
	}

	if tgv.Volley.Price > 0 {
		text += fmt.Sprintf("\n💰 %d ₽", tgv.Volley.Price)
	}

	if tgv.Volley.CourtCount > 0 {
		text += fmt.Sprintf("\n*Корты:* %d", tgv.Volley.CourtCount)
	}
	if tgv.Volley.MaxPlayers > 0 {
		text += fmt.Sprintf("\n*Игроков:* %d", tgv.Volley.MaxPlayers)
	}
	text += tgv.GetPlayersText()
	return
}

func (tgv *TelegramView) GetPlayersText() (text string) {
	count := 1
	over := false
	for _, pl := range tgv.Volley.Players {
		if pl.Count == 0 {
			continue
		}
		pvw := person.NewTelegramViewRu(pl.Person)
		text += fmt.Sprintf("\n%d. %s", count, pvw.String())
		if !pl.ArriveTime.IsZero() {
			text += fmt.Sprintf(" (%s)", pl.ArriveTime.Format("15:04"))
		}
		count++
		if !over && count > tgv.Volley.MaxPlayers {
			over = true
			text += "\n\n*Резерв:*"
			count = 1
		}
		for i := 1; i < pl.Count; i++ {
			text += fmt.Sprintf("\n%d. %s+%d", count, pl.String(), i)
			count++
			if !over && count > tgv.Volley.MaxPlayers {
				over = true
				text += "\n\n*Резерв:*"
				count = 1
			}
		}
	}
	if !over && tgv.Volley.MaxPlayers-count > 3 {
		text += fmt.Sprintf("\n%d.\n.\n.\n%d.", count, tgv.Volley.MaxPlayers)
	} else if !over && tgv.Volley.MaxPlayers > 0 {
		for i := count; i <= tgv.Volley.MaxPlayers; i++ {
			text += fmt.Sprintf("\n%d.", i)
		}
	}
	if tgv.Reserve.Description != "" {
		text += "\n\n" + tgv.Reserve.Description
	}
	return
}

func (tgv *TelegramView) GetTimeText() (text string) {
	if !tgv.Reserve.StartTime.IsZero() {
		text += tgv.Reserve.StartTime.Format("15:04")
		if !tgv.Reserve.GetEndTime().IsZero() {
			text += fmt.Sprintf("-%s", tgv.Reserve.EndTime.Format("15:04"))
		}
	}
	return
}
