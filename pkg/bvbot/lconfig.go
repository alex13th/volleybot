package bvbot

import (
	"volleybot/pkg/telegram"
)

type ConfigStateProvider struct {
	BaseStateProvider
	Resources ConfigResources
}

func (p ConfigStateProvider) GetMR() *telegram.MessageRequest {
	cfgview := NewConfigTelegramViewRu(p.GetLocationConfig())
	txt := cfgview.GetText()

	if p.kh == nil {
		p.kh = p.GetKeyboardHelper()
	}
	kbd := p.kh.GetKeyboard()

	return p.CreateMR(p.State.ChatId, txt, p.Resources.ParseMode, kbd)
}

func (p ConfigStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	var sreq telegram.StateRequest
	sreq.State = p.State
	sreq.Request = p.GetEditMR(p.GetMR())

	return append(reqlist, sreq)
}

func (p ConfigStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 1
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgcourts", Text: res.Courts.CourtBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgprice", Text: res.Price.PriceBtn})
	}
	return &ah
}
