package bvbot

import (
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"
)

type BaseStateProvider struct {
	reserve    volley.Volley
	kh         telegram.KeyboardHelper
	BackState  telegram.State
	Message    telegram.Message
	Person     person.Person
	Repository volley.Repository
	Location   location.Location
	State      telegram.State
	Text       string
}

func NewBaseStateProvider(state telegram.State, msg telegram.Message, person person.Person, loc location.Location, repository volley.Repository, text string) (sp BaseStateProvider, err error) {
	sp = BaseStateProvider{State: state, Message: msg, Person: person, Location: loc, Repository: repository, Text: text}
	if repository != nil && state.Data != "" {
		id, err := volley.Volley{}.IdFromBase64(state.Data)
		if err != nil {
			return sp, err
		}
		sp.reserve, _ = sp.Repository.Get(id)
	}
	return
}

func (p BaseStateProvider) GetBaseKeyboardHelper(txt string) (kh telegram.BaseKeyboardHelper) {
	kh.Text = txt
	kh.State = p.State
	if p.Message.Chat.Id > 0 {
		kh.BackData = p.BackState.String()
	}
	return
}

func (p BaseStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Updated {
		err = p.Repository.Update(p.reserve)
	}
	st = p.State
	st.Data = p.reserve.Base64Id()
	st.State = p.State.Action
	return
}

func (p BaseStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action != p.State.State {
		return
	}
	if p.Message.Chat.Id < 0 {
		req := telegram.DeleteMessageRequest{ChatId: p.Message.Chat.Id, MessageId: p.Message.MessageId}
		rlist = append(rlist, telegram.StateRequest{Request: &req})
		rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetMR()})

	} else {
		rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetEditMR()})
	}
	return
}

func (p *BaseStateProvider) GetEditMR() (mer *telegram.EditMessageTextRequest) {
	mr := p.GetMR()
	mer = &telegram.EditMessageTextRequest{MessageId: p.Message.MessageId, ChatId: mr.ChatId, Text: mr.Text, ParseMode: mr.ParseMode,
		ReplyMarkup: mr.ReplyMarkup}
	return
}

func (p *BaseStateProvider) GetMR() (mr *telegram.MessageRequest) {
	cid := p.Message.Chat.Id
	rview := volley.NewTelegramViewRu(p.reserve)
	mtxt := rview.GetText()

	var kbd telegram.InlineKeyboardMarkup
	if p.kh != nil {
		kbd.InlineKeyboard = p.kh.GetKeyboard()
		if kbdText := p.kh.GetText(); kbdText != "" {
			mtxt += "\n" + kbdText
		}
	}

	if cid < 0 {
		mtxt += "\n\n"
		mtxt += p.Text
	}

	if len(kbd.InlineKeyboard) > 0 {
		return &telegram.MessageRequest{ChatId: cid, Text: mtxt, ParseMode: rview.ParseMode,
			ReplyMarkup: kbd, DisableNotification: cid < 0}
	}
	return &telegram.MessageRequest{ChatId: cid, Text: mtxt, ParseMode: rview.ParseMode,
		DisableNotification: cid < 0}
}

func (p *BaseStateProvider) NotifyPlayers(action string) (reqlist []telegram.StateRequest) {
	for _, pl := range p.reserve.Players {
		if pl.Person.TelegramId != p.Person.TelegramId {
			if param, ok := pl.Settings[action]; ok && param == "on" {
				mr := p.GetMR()
				mr.ChatId = pl.TelegramId
				reqlist = append(reqlist, telegram.StateRequest{Request: mr})
			}
		}
	}
	return
}

type MessageProcessorResources struct {
	Actions      ActionsResources
	Activity     AcivityResources
	Courts       CourtsResources
	Cancel       CancelResources
	Description  DescResources
	Join         JoinResources
	Level        LevelResources
	List         ListResources
	Main         MainResources
	MaxPlayer    MaxPlayersResources
	RemovePlayer RemovePlayerResources
	Price        PriceResources
	Settings     SettingsResources
	Sets         SetsResources
	Show         ShowResources
	BackBtn      string
	DescMessage  string
}
