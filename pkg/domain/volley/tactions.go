package volley

import (
	"strconv"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type ActionsResources struct {
	BackBtn         string
	CancelBtn       string
	CopyBtn         string
	CopyDoneMessage string
	PublishBtn      string
	RemovePlayerBtn string
}

func NewActionsResourcesRu() (r ActionsResources) {
	r.BackBtn = "Назад"
	r.CancelBtn = "💥Отменить"
	r.CopyBtn = "🫂 Копировать"
	r.CopyDoneMessage = "Копия сделана! 👆"
	r.PublishBtn = "Опубликовать"
	r.RemovePlayerBtn = "Удалить игрока"
	return
}

type ActionsStateProvider struct {
	BaseStateProvider
	Resources     ActionsResources
	ShowResources ShowResources
}

func (p ActionsStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action == "pub" {
		show_p := ShowStateProvider{BaseStateProvider: p.BaseStateProvider, Resources: p.ShowResources}
		show_p.Message.Chat = &telegram.Chat{Id: p.Location.ChatId}
		show_p.State.State = "show"
		show_p.State.Action = "show"
		return show_p.GetRequests()
	}
	if p.State.Action == "copy" {
		req := &telegram.MessageRequest{ChatId: p.Message.Chat.Id, Text: p.Resources.CopyDoneMessage}
		return append(rlist, telegram.StateRequest{State: p.State, Request: req})
	}
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p ActionsStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	msg := p.Message
	kh := telegram.ActionsKeyboardHelper{}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	kh.Actions = []telegram.ActionButton{}
	if p.reserve.Canceled {
		return &kh
	}
	kh.Columns = 2
	if msg.Chat.Id == p.Person.TelegramId {
		if p.reserve.Person.TelegramId == p.Person.TelegramId || p.Person.CheckLocationRole(p.reserve.Location, "admin") {
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "cancel", Text: res.CancelBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "copy", Text: res.CopyBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "pub", Text: res.PublishBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "rmpl", Text: res.RemovePlayerBtn})
		}
	}
	return &kh
}

func (p ActionsStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Action == "copy" {
		p.State.Action = "show"
		if p.reserve, err = p.Repository.Add(p.reserve.Copy()); err != nil {
			return
		}
	}
	if p.State.Action == "pub" {
		return p.BackState, nil
	}
	return p.BaseStateProvider.Proceed()
}

type CancelResources struct {
	BackBtn    string
	Text       string
	ConfirmBtn string
	AbortBtn   string
}

func NewCancelResourcesRu() (r CancelResources) {
	r.BackBtn = "Передумал"
	r.ConfirmBtn = "🧨 Уверен"
	r.Text = "\n🧨*ВНИМАНИЕ!!!*🧨\nИгра будет отменена для всех участников. Если есть желание только выписаться, лучше воспользоваться кнопкой \"Не буду\""
	return
}

type CancelStateProvider struct {
	BaseStateProvider
	Resources     CancelResources
	ShowResources ShowResources
}

func (p CancelStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p CancelStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	msg := p.Message
	kh := telegram.ActionsKeyboardHelper{}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(p.Resources.Text)
	kh.Actions = []telegram.ActionButton{}
	if p.reserve.Canceled {
		return &kh
	}
	kh.Columns = 2
	if msg.Chat.Id == p.Person.TelegramId {
		if p.reserve.Person.TelegramId == p.Person.TelegramId || p.Person.CheckLocationRole(p.reserve.Location, "admin") {
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "confirm", Text: res.ConfirmBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "leave", Text: p.ShowResources.JoinLeaveBtn})
		}
	}
	return &kh
}

func (p CancelStateProvider) Proceed() (telegram.State, error) {
	if p.State.Action == "confirm" {
		p.reserve.Canceled = true
		p.State.Action = "show"
		p.State.Updated = true
	}
	if p.State.Action == "leave" {
		rpl := p.reserve.GetPlayer(p.Person.Id)
		if rpl.Id == uuid.Nil {
			rpl = person.Player{Person: p.Person}
		}
		rpl.Count = 0
		p.reserve.JoinPlayer(person.Player{Person: p.Person, Count: 0})
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
}

type RemovePlayerResources struct {
	BackBtn         string
	Message         string
	RemovePlayerBtn string
}

func RemovePlayerResourcesRu() (r RemovePlayerResources) {
	r.BackBtn = "Назад"
	r.RemovePlayerBtn = "Удалить игрока"
	return
}

type RemovePlayerStateProvider struct {
	BaseStateProvider
	Resources RemovePlayerResources
}

func (p RemovePlayerStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p RemovePlayerStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	msg := p.Message
	if msg.Chat.Id == p.Person.TelegramId {
		pllist := []telegram.EnumItem{}
		for _, pl := range p.reserve.Players {
			pllist = append(pllist, telegram.EnumItem{Id: strconv.Itoa(pl.TelegramId), Item: pl.String()})
		}
		kh := telegram.NewEnumKeyboardHelper(pllist)
		kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(p.Resources.Message)
		return &kh
	}
	return nil
}

func (p RemovePlayerStateProvider) Proceed() (telegram.State, error) {
	if p.State.Action == "set" {
		kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
		ptid, _ := strconv.Atoi(kh.Value)
		rpl := p.reserve.GetPlayerByTelegramId(ptid)
		rpl.Count = 0
		p.reserve.JoinPlayer(rpl)
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
}
