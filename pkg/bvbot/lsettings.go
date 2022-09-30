package bvbot

import (
	"volleybot/pkg/telegram"
)

type ConfigStateProvider struct {
	BaseStateProvider
	Config    Config
	Resources ConfigResources
}

func (p ConfigStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Updated {
		p.ConfigRepository.Update(p.Location, p.name, p.Config)
		p.State.Updated = false
	}
	return p.BaseStateProvider.Proceed()
}

func (p ConfigStateProvider) GetMR() *telegram.MessageRequest {
	cfgview := NewConfigTelegramViewRu(p.Config)
	txt := cfgview.GetText()

	var kbd telegram.InlineKeyboardMarkup
	kbd.InlineKeyboard = p.GetKeyboardHelper().GetKeyboard()

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
			Action: "lcourts", Text: res.Courts.CourtBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "lprice", Text: res.Price.PriceBtn})
	}
	return &ah
}
