package bvbot

import (
	"fmt"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/telegram"
)

type ProfileResources struct {
	LevelBtn    string
	NotifiesBtn string
	ParseMode   string
	SexBtn      string
	Text        string
}

func ProfileResourcesRu() (r ProfileResources) {
	r.LevelBtn = "Уровень"
	r.NotifiesBtn = "Оповещения"
	r.ParseMode = "Markdown"
	r.SexBtn = "Пол"
	r.Text = ""
	return
}

type ProfileStateProvider struct {
	BaseStateProvider
	Resources ProfileResources
}

func (p ProfileStateProvider) GetMR() *telegram.MessageRequest {
	pview := person.NewTelegramViewRu(p.Person)
	psetview := person.NewTelegramSettingsViewRu(p.Person)
	txt := fmt.Sprintf("%s\n\n%s\n%s", pview.GetText(), psetview.GetText(), p.Resources.Text)

	var kbd telegram.InlineKeyboardMarkup
	kh := p.GetKeyboardHelper()
	kbd.InlineKeyboard = kh.GetKeyboard()

	return p.CreateMR(p.State.ChatId, txt, p.Resources.ParseMode, kbd)
}

func (p ProfileStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	var sreq telegram.StateRequest
	sreq.State = p.State
	sreq.Request = p.GetEditMR(p.GetMR())
	return append(reqlist, sreq)
}

func (p ProfileStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "level", Text: res.LevelBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "sex", Text: res.SexBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "notifies", Text: res.NotifiesBtn})
	}
	return &ah
}

func (p ProfileStateProvider) Proceed() (st telegram.State, err error) {
	return p.BaseStateProvider.Proceed()
}
