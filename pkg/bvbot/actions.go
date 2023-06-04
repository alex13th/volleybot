package bvbot

import (
	"strconv"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ActionsStateProvider struct {
	BaseStateProvider
	Resources     ActionsResources
	ShowResources ShowResources
}

func (p ActionsStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action == "pub" {
		show_p := ShowStateProvider{BaseStateProvider: p.BaseStateProvider, Resources: p.ShowResources}
		show_p.State.ChatId = p.Location.ChatId
		show_p.State.State = "show"
		show_p.State.Action = "show"
		return show_p.GetRequests()
	}
	if p.State.Action == "copy" {
		req := &telegram.MessageRequest{ChatId: p.State.ChatId, Text: p.Resources.CopyDoneMessage}
		return append(rlist, telegram.StateRequest{State: p.State, Request: req})
	}
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p ActionsStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.ActionsKeyboardHelper{}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	kh.Actions = []telegram.ActionButton{}
	if p.reserve.Canceled {
		return &kh
	}
	kh.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
		if p.reserve.Person.TelegramId == p.Person.TelegramId || p.Person.CheckLocationRole(p.reserve.Location, "admin") {
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "cancel", Text: res.CancelBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "copy", Text: res.CopyBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "paid", Text: res.PaidBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "rmpl", Text: res.RemovePlayerBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "pub", Text: res.PublishBtn})
			kh.Actions = append(kh.Actions, telegram.ActionButton{
				Action: "send", Text: res.SendBtn})
		}
	}
	return &kh
}

func (p ActionsStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Action == "copy" {
		if p.reserve, err = p.Repository.Add(p.reserve.Copy()); err != nil {
			p.State.Action = "show"
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ActionsStateProvider",
				"provider": p,
				"state":    p.State,
				"error":    err,
			}).Error("error copy reserve")
			return
		}
	}
	if p.State.Action == "pub" {
		return p.BackState, nil
	}
	return p.BaseStateProvider.Proceed()
}

type CancelStateProvider struct {
	BaseStateProvider
	Resources     CancelResources
	ShowResources ShowResources
}

func (p CancelStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	p.kh = p.GetKeyboardHelper()
	reqlist = p.BaseStateProvider.GetRequests()
	if p.State.Action == "confirm" {
		p.reserve.Canceled = true
		p.kh = nil
		reqlist = append(reqlist, p.NotifyPlayers("notify_cancel")...)
	}
	return
}

func (p CancelStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.ActionsKeyboardHelper{}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(p.Resources.Text)
	kh.Actions = []telegram.ActionButton{}
	if p.reserve.Canceled {
		return &kh
	}
	kh.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
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
		mb := p.reserve.GetMember(p.Person.Id)
		if mb.Id == uuid.Nil {
			mb = volley.Member{Player: volley.NewPlayer(p.Person)}
		}
		mb.Count = 0
		p.reserve.JoinPlayer(volley.Member{Player: mb.Player, Count: 0})
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
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
	if p.State.ChatId == p.Person.TelegramId {
		pllist := []telegram.EnumItem{}
		for _, mb := range p.reserve.Members {
			pllist = append(pllist, telegram.EnumItem{Id: strconv.Itoa(mb.TelegramId), Item: mb.String()})
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
		ptid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "RemovePlayerStateProvider",
				"provider": p,
				"state":    p.State,
				"error":    err,
			}).Error("can't parse player telegram id: " + kh.Value)
		}
		rpl := p.reserve.GetMemberByTelegramId(ptid)
		rpl.Count = 0
		p.reserve.JoinPlayer(rpl)
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
}

type PaidPlayerStateProvider struct {
	BaseStateProvider
	Resources RemovePlayerResources
}

func (p PaidPlayerStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p PaidPlayerStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	if p.State.ChatId == p.Person.TelegramId {
		pllist := []telegram.EnumItem{}
		for _, mb := range p.reserve.Members {
			pllist = append(pllist, telegram.EnumItem{Id: strconv.Itoa(mb.TelegramId), Item: mb.String()})
		}
		kh := telegram.NewEnumKeyboardHelper(pllist)
		kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(p.Resources.Message)
		return &kh
	}
	return nil
}

func (p *PaidPlayerStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Action == "set" {
		kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
		ptid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "PaidPlayerStateProvider",
				"provider": p,
				"state":    p.State,
				"error":    err,
			}).Error("can't parse player telegram id: " + kh.Value)
		}
		rpl := p.reserve.GetMemberByTelegramId(ptid)
		rpl.SetPaid(!rpl.GetPaid())
		p.reserve.JoinPlayer(rpl)
		p.State.Action = p.State.State
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
}

type SendStateProvider struct {
	BaseStateProvider
	Resources SendResources
}

func (p SendStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action == "send" {
		kh := telegram.SendKeyboardHelper{Text: p.Resources.SendBtn, RequestId: p.Message.MessageId}
		req := telegram.MessageRequest{ChatId: p.State.ChatId, Text: p.Resources.Message, ReplyMarkup: kh.GetKeyboard()}
		p.State.MessageId = -1
		p.State.Action = "done"
		return append(rlist, telegram.StateRequest{State: p.State, Request: &req})
	}
	p.State.ChatId = p.Message.From.Id
	req := telegram.MessageRequest{ChatId: p.State.ChatId, Text: "LASLALA", ReplyMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true}}
	return append(rlist, telegram.StateRequest{State: p.State, Clear: true, Request: req})
}

func (p *SendStateProvider) Proceed() (telegram.State, error) {
	if p.State.State == "send" && p.Message.ChatShared != nil {
		p.State.ChatId = p.Message.ChatShared.ChatId
		p.State.Action = "show"
	}
	return p.BaseStateProvider.Proceed()
}
