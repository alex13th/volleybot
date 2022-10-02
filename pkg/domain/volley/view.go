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
	for _, pl := range tgv.Volley.Members {
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
		text += fmt.Sprintf("\n💪*Уровень*: %s", PlayerLevel(tgv.Volley.MinLevel))
	}
	if tgv.Volley.NetType > 0 {
		text += fmt.Sprintf("\n*Сетка*: %s", NetType(tgv.Volley.NetType))
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
	text += tgv.GetMembersText()
	return
}

func (tgv *TelegramView) GetMembersText() (text string) {
	count := 1
	over := false
	for _, mb := range tgv.Volley.Members {
		if mb.Count == 0 {
			continue
		}
		pvw := NewPlayerTelegramView(mb.Player)
		text += fmt.Sprintf("\n%d. %s", count, pvw.String())
		if !mb.ArriveTime.IsZero() {
			text += fmt.Sprintf(" (%s)", mb.ArriveTime.Format("15:04"))
		}
		if mb.Paid {
			text += " 💴"
		}
		count++
		if !over && count > tgv.Volley.MaxPlayers {
			over = true
			text += "\n\n*Резерв:*"
			count = 1
		}
		for i := 1; i < mb.Count; i++ {
			text += fmt.Sprintf("\n%d. %s+%d", count, mb.String(), i)
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

type PlayerTelegramView struct {
	Player
	ParseMode string
}

func NewPlayerTelegramView(p Player) PlayerTelegramView {
	return PlayerTelegramView{Player: p, ParseMode: "Markdown"}
}

func (tgv *PlayerTelegramView) String() (text string) {
	pv := person.NewTelegramViewRu(tgv.Person)
	text = PlayerLevel(tgv.Level).Emoji()
	text += pv.String()
	return
}

func (tgv PlayerTelegramView) GetText() (text string) {
	pv := person.NewTelegramViewRu(tgv.Person)
	text = pv.GetText()
	text += fmt.Sprintf("\n*Уровень*: %s", tgv.GetLevelText())
	return
}

func (tgv *PlayerTelegramView) GetLevelText() (text string) {
	if tgv.Level > 0 {
		text = PlayerLevel(tgv.Level).Emoji() + " "
	}
	text += PlayerLevel(tgv.Level).String()
	return
}
