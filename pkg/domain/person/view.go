package person

type PersonView interface {
	GetText() (text string)
}

type TelegramView struct {
	Person Person
}

func (tgv *TelegramView) GetText() (text string) {
	text = tgv.Person.GetDisplayname()
	return
}
