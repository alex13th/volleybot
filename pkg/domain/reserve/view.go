package reserve

import (
	"fmt"
	"sort"

	"github.com/goodsign/monday"
	"github.com/google/uuid"
)

type ReserveView interface {
	GetText() (text string)
}

func NewTelegramViewRu(res Reserve) TelegramView {
	return TelegramView{
		Reserve:     res,
		CancelLabel: "🔥*ОТМЕНА*🔥",
		DateLabel:   "📆",
		TimeLabel:   "⏰",
		Locale:      monday.LocaleRuRU,
		ParseMode:   "Markdown",
	}
}

type TelegramView struct {
	Reserve     Reserve
	CancelLabel string
	DateLabel   string
	TimeLabel   string
	Locale      monday.Locale
	ParseMode   string
}

func (tgv *TelegramView) String() string {
	plcount := 0
	for _, pl := range tgv.Reserve.Players {
		plcount += pl.Count
	}
	return fmt.Sprintf("%s %s (%d/%d)",
		monday.Format(tgv.Reserve.StartTime, "Mon, 02.01", tgv.Locale),
		tgv.GetTimeText(), plcount, tgv.Reserve.MaxPlayers)
}

func (tgv *TelegramView) GetText() (text string) {
	if tgv.Reserve.Canceled {
		text = tgv.CancelLabel + "\n"
	}
	var uname string
	if tgv.Reserve.Person.TelegramId != 0 {
		uname = "[%s](tg://user?id=%d)"
		uname = fmt.Sprintf(uname, tgv.Reserve.Person.GetDisplayname(), tgv.Reserve.Person.TelegramId)
	} else {
		uname = fmt.Sprintf("*%s*", tgv.Reserve.Person.GetDisplayname())
	}
	text += fmt.Sprintf("%s\n%s %s\n%s %s", uname,
		tgv.DateLabel, monday.Format(tgv.Reserve.StartTime, "Monday, 02.01.2006", tgv.Locale),
		tgv.TimeLabel, tgv.GetTimeText())
	if tgv.Reserve.MinLevel > 0 {
		text += fmt.Sprintf("\n💪*Уровень*: %s", PlayerLevel(tgv.Reserve.MinLevel))
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
	count := 0
	keys := make([]string, 0, len(tgv.Reserve.Players))
	for k := range tgv.Reserve.Players {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)

	for _, k := range keys {
		var uname string
		id, _ := uuid.Parse(k)
		pl := tgv.Reserve.Players[id]
		if pl.Person.TelegramId != 0 {
			uname = "[%s](tg://user?id=%d)"
			uname = fmt.Sprintf(uname, pl.Person.GetDisplayname(), pl.Person.TelegramId)
		} else {
			uname = pl.Person.GetDisplayname()
		}
		text += fmt.Sprintf("\n%d. %s", count+1, uname)
		for i := 1; i < pl.Count; i++ {
			text += fmt.Sprintf("\n%d. %s+%d", count+i+1, uname, i)
		}
		count += pl.Count
	}
	if tgv.Reserve.MaxPlayers-count-1 > 3 {
		text += fmt.Sprintf("\n%d.\n.\n.\n%d.", count+1, tgv.Reserve.MaxPlayers)
	} else if tgv.Reserve.MaxPlayers > 0 {
		for i := count + 1; i <= tgv.Reserve.MaxPlayers; i++ {
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
