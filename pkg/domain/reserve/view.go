package reserve

import (
	"fmt"
	"volleybot/pkg/domain/person"

	"github.com/goodsign/monday"
)

type ReserveView interface {
	GetText() (text string)
}

func NewTelegramViewRu(res Reserve) TelegramView {
	return TelegramView{
		Reserve:         res,
		CancelLabel:     "🔥 *ОТМЕНА* 🔥",
		GameLabel:       "🏐 *СВОБОДНЫЕ ИГРЫ* 🏐",
		TrainingLabel:   "‼️ *ТРЕНИРОВКА* ‼️",
		TournamentLabel: "💥🔥 *ТУРНИР* 🔥💥",
		TennisLabel:     "🎾 *ПЛЯЖНЫЙ ТЕННИС* 🎾",
		DateLabel:       "📆",
		TimeLabel:       "⏰",
		Locale:          monday.LocaleRuRU,
		ParseMode:       "Markdown",
	}
}

type TelegramView struct {
	Reserve         Reserve
	CancelLabel     string
	GameLabel       string
	TrainingLabel   string
	TournamentLabel string
	TennisLabel     string
	DateLabel       string
	TimeLabel       string
	Locale          monday.Locale
	ParseMode       string
}

func (tgv *TelegramView) String() string {
	plcount := 0
	for _, pl := range tgv.Reserve.Players {
		plcount += pl.Count
	}
	return fmt.Sprintf("%s %s %s (%d/%d)", tgv.Reserve.Activity.Emoji(),
		monday.Format(tgv.Reserve.StartTime, "Mon, 02.01", tgv.Locale),
		tgv.GetTimeText(), plcount, tgv.Reserve.MaxPlayers)
}

func (tgv *TelegramView) GetText() (text string) {
	if tgv.Reserve.Canceled {
		text = tgv.CancelLabel + "\n\n"
	} else if tgv.Reserve.Activity == 10 {
		text = tgv.TrainingLabel + "\n\n"
	} else if tgv.Reserve.Activity == 20 {
		text = tgv.TournamentLabel + "\n\n"
	} else if tgv.Reserve.Activity == 30 {
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
	if tgv.Reserve.MinLevel > 0 {
		text += fmt.Sprintf("\n💪*Уровень*: %s", person.PlayerLevel(tgv.Reserve.MinLevel))
	}

	if tgv.Reserve.Price > 0 {
		text += fmt.Sprintf("\n💰 %d ₽", tgv.Reserve.Price)
	}

	if tgv.Reserve.CourtCount > 0 {
		text += fmt.Sprintf("\n*Корты:* %d", tgv.Reserve.CourtCount)
	}
	if tgv.Reserve.MaxPlayers > 0 {
		text += fmt.Sprintf("\n*Игроков:* %d", tgv.Reserve.MaxPlayers)
	}
	text += tgv.GetPlayersText()
	return
}

func (tgv *TelegramView) GetPlayersText() (text string) {
	count := 1
	over := false
	for _, pl := range tgv.Reserve.Players {
		pvw := person.NewTelegramViewRu(pl.Person)
		text += fmt.Sprintf("\n%d. %s", count, pvw.String())
		if !pl.ArriveTime.IsZero() {
			text += fmt.Sprintf(" (%s)", pl.ArriveTime.Format("15:04"))
		}
		count++
		if !over && count > tgv.Reserve.MaxPlayers {
			over = true
			text += "\n\n*Резерв:*"
			count = 1
		}
		for i := 1; i < pl.Count; i++ {
			text += fmt.Sprintf("\n%d. %s+%d", count, pl.String(), i)
			count++
			if !over && count > tgv.Reserve.MaxPlayers {
				over = true
				text += "\n\n*Резерв:*"
				count = 1
			}
		}
	}
	if !over && tgv.Reserve.MaxPlayers-count > 3 {
		text += fmt.Sprintf("\n%d.\n.\n.\n%d.", count, tgv.Reserve.MaxPlayers)
	} else if !over && tgv.Reserve.MaxPlayers > 0 {
		for i := count; i <= tgv.Reserve.MaxPlayers; i++ {
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
