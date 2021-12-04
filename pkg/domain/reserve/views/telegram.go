package views

import (
	"fmt"
	"volleybot/pkg/domain/reserve"

	"github.com/goodsign/monday"
)

func NewTelegramViewRu(res reserve.Reserve) TelegramView {
	return TelegramView{
		Reserve:   res,
		DateLabel: "📆",
		TimeLabel: "⏰",
		Locale:    monday.LocaleRuRU,
		ParseMode: "Markdown",
	}
}

type TelegramView struct {
	Reserve   reserve.Reserve
	DateLabel string
	TimeLabel string
	Locale    monday.Locale
	ParseMode string
}

func (tgv *TelegramView) GetText() (text string) {
	tt := "*%s*\n%s"
	tgv.Reserve.Person.GetDisplayname()
	text = fmt.Sprintf(
		tt, tgv.Reserve.Person.GetDisplayname(),
		tgv.GetTimeText())
	return
}

func (tgv *TelegramView) GetTimeText() (text string) {
	tt := "%s %s\n%s %s-%s"
	text = fmt.Sprintf(
		tt,
		tgv.DateLabel, monday.Format(tgv.Reserve.StartTime, "Monday, 02.01.2006", tgv.Locale),
		tgv.TimeLabel, tgv.Reserve.StartTime.Format("15:04"), tgv.Reserve.EndTime.Format("15:04"))
	return
}
