package person

import "fmt"

type PersonView interface {
	GetText() (text string)
}

type TelegramView struct {
	Person    Person
	ParseMode string
}

func NewTelegramViewRu(p Person) TelegramView {
	return TelegramView{
		Person:    p,
		ParseMode: "Markdown",
	}
}

func (tgv *TelegramView) GetText() (text string) {
	text = fmt.Sprintf("*Имя*: %s", tgv.Person.Firstname)
	text += fmt.Sprintf("\n*Фамилия*: %s", tgv.Person.Lastname)
	text += fmt.Sprintf("\n*Полное имя*: %s", tgv.Person.GetDisplayname())
	text += fmt.Sprintf("\n*Пол*: %s", tgv.GetSexText())
	text += fmt.Sprintf("\n*Уровень*: %s", tgv.GetLevelText())
	return
}

func (tgv *TelegramView) GetLevelText() (text string) {
	if tgv.Person.Level > 0 {
		text = PlayerLevel(tgv.Person.Level).Emoji() + " "
	}
	text += PlayerLevel(tgv.Person.Level).String()
	return
}

func (tgv *TelegramView) GetSexText() (text string) {
	if tgv.Person.Sex != 0 {
		text = Sex(tgv.Person.Sex).Emoji() + " " + tgv.Person.Sex.String()
	}
	return
}

func (tgv *TelegramView) String() (text string) {
	text = tgv.Person.GetDisplayname()
	text = fmt.Sprintf("%s %s", Sex.Emoji(tgv.Person.Sex), text)
	text = fmt.Sprintf("%s%s", PlayerLevel(tgv.Person.Level).Emoji(), text)

	if tgv.Person.TelegramId != 0 {
		text = fmt.Sprintf("[%s](tg://user?id=%d)", text, tgv.Person.TelegramId)
	}

	return
}
