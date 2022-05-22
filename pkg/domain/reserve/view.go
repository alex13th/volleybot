package reserve

import (
	"fmt"

	"github.com/goodsign/monday"
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

func (tgv *TelegramView) GetText() (text string) {
	tt := ""
	if tgv.Reserve.Canceled {
		tt = tgv.CancelLabel + "\n"
	}
	tt += "*%s*\n%s"
	text = fmt.Sprintf(
		tt, tgv.Reserve.Person.GetDisplayname(),
		tgv.GetTimeText())

	if tgv.Reserve.MinLevel > 0 {
		text += fmt.Sprintf("\n💪*Уровень*: %s", PlayerLevel(tgv.Reserve.MinLevel))
	}

	if tgv.Reserve.Price > 0 {
		text += fmt.Sprintf("\n💳 %d", tgv.Reserve.Price)
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
	for _, pl := range tgv.Reserve.Players {
		if pl.Count > 1 {
			text += fmt.Sprintf("\n%d. %s+%d", count+1, pl.Person.GetDisplayname(), pl.Count-1)
		} else {
			text += fmt.Sprintf("\n%d. %s", count+1, pl.Person.GetDisplayname())
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
	return
}

func (tgv *TelegramView) GetTimeText() (text string) {
	text = fmt.Sprintf("%s %s", tgv.DateLabel,
		monday.Format(tgv.Reserve.StartTime, "Monday, 02.01.2006", tgv.Locale))
	if tgv.Reserve.StartTime.Hour() > 0 {
		text += fmt.Sprintf("\n%s %s", tgv.TimeLabel,
			tgv.Reserve.StartTime.Format("15:04"))
		if !tgv.Reserve.EndTime.IsZero() {
			text += fmt.Sprintf("-%s", tgv.Reserve.EndTime.Format("15:04"))
		}
	}
	return
}
