package views

import "volleybot/pkg/domain/person"

type TelegramView struct {
	Person person.Person
}

func (tgv *TelegramView) GetText() (text string) {
	text = tgv.Person.GetDisplayname()
	return
}
