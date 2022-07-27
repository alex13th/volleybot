package reserve

import (
	"fmt"

	"github.com/goodsign/monday"
)

type TelegramViewResources struct {
	CancelLabel string
	DateLabel   string
	TimeLabel   string
	Locale      monday.Locale
	ParseMode   string
}

func NewTelegramResourcesRu() TelegramViewResources {
	return TelegramViewResources{
		CancelLabel: "🔥 *ОТМЕНА* 🔥",
		DateLabel:   "📆",
		TimeLabel:   "⏰",
		Locale:      monday.LocaleRuRU,
		ParseMode:   "Markdown",
	}
}

type ReserveView interface {
	GetText() (text string)
}

func NewTelegramViewRu(res Reserve) TelegramView {
	return TelegramView{Reserve: res, TelegramViewResources: NewTelegramResourcesRu()}
}

type TelegramView struct {
	Reserve
	TelegramViewResources
}

func (tgv *TelegramView) String() string {
	return fmt.Sprintf("%s %s",
		monday.Format(tgv.Reserve.StartTime, "Mon, 02.01", tgv.Locale),
		tgv.GetTimeText())
}

func (tgv *TelegramView) GetText() (text string) {
	if tgv.Reserve.Canceled {
		text = tgv.CancelLabel + "\n\n"
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

	if tgv.Reserve.Price > 0 {
		text += fmt.Sprintf("\n💰 %d ₽", tgv.Reserve.Price)
	}

	text += tgv.GetPlayersText()
	return
}

func (tgv *TelegramView) GetPlayersText() (text string) {
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
