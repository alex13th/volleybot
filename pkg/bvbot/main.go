package bvbot

import (
	"fmt"
	"time"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type MainStateProvider struct {
	BaseStateProvider
	Resources MainResources
}

func (p MainStateProvider) GetMR() (mr *telegram.MessageRequest) {
	txt := ""
	txt += p.Resources.Text

	kh := p.GetKeyboardHelper()
	kbd := kh.GetKeyboard()

	return p.CreateMR(p.State.ChatId, txt, p.Resources.ParseMode, kbd)
}

func (p MainStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.ChatId < 0 {
		return
	}
	var sreq telegram.StateRequest
	if p.State.Action == "start" {
		sreq.State = p.State
		sreq.Request = p.GetMR()
		return append(rlist, sreq)
	}
	if p.State.Action == "main" {
		sreq.State = p.State
		sreq.Request = p.GetEditMR(p.GetMR())
		return append(rlist, sreq)
	}
	return
}

func (p MainStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 1
	if p.State.ChatId == p.Person.TelegramId {
		if p.Person.CheckLocationRole(p.Location, "admin") || p.Person.CheckLocationRole(p.Location, "order") {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "order", Text: res.NewReserveBtn})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "today", Text: res.TodayBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "listd", Text: res.ListDateBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "profile", Text: res.ProfileBtn})
	}
	if p.Person.CheckLocationRole(p.Location, "admin") {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "config", Text: res.ConfigBtn})
	}
	return &ah
}

func (p MainStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Action == "order" {
		p.reserve = p.NewReserve()
		p.reserve.Location = p.Location
		p.reserve, err = p.Repository.Add(p.reserve)
		p.State.Data = p.reserve.Base64Id()
		p.State.Action = "show"
	} else if p.State.Action == "listd" {
		return p.BaseStateProvider.Proceed()
	} else if p.State.Action == "profile" {
		return p.BaseStateProvider.Proceed()
	} else if p.State.Action == "config" {
		return p.BaseStateProvider.Proceed()
	} else if p.State.Action == "today" {
		p.State.State = "listd"
		p.State.Action = "set"
		p.State.Value = time.Now().Format("2006-01-02")
		return p.State, err
	} else {
		p.State.Action = ""
		return p.State, err
	}
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "bvbot",
			"function": "Proceed",
			"struct":   "MainStateProvider",
			"error":    err,
		}).Error("main state proceed error")
		return
	}
	return p.BaseStateProvider.Proceed()
}

func (p MainStateProvider) NewReserve() (r volley.Volley) {
	currTime := time.Now()
	stime := time.Date(currTime.Year(), currTime.Month(), currTime.Day(),
		currTime.Hour()+1, 0, 0, 0, currTime.Location())
	etime := stime.Add(time.Duration(time.Hour))

	r = volley.NewVolley(p.Person, stime, etime)
	return
}

type ListdStateProvider struct {
	reserves []volley.Volley
	BaseStateProvider
	Resources ListResources
}

func (p ListdStateProvider) GetListText() (txt string) {
	if len(p.reserves) > 0 {
		txt += p.Resources.ListCaption + "\n\n"
		for i, res := range p.reserves {
			tgv := volley.NewTelegramViewRu(res)
			txt += fmt.Sprintf("%v. %s\n", i+1, tgv.String())
		}
		txt += "\n"
	} else {
		txt += p.Resources.NoReservesMessage
	}
	return
}

func (p ListdStateProvider) GetMR() (mr *telegram.MessageRequest) {
	txt := ""
	if p.State.Action == "list" {
		txt += p.GetListText()
	}
	kh := p.GetKeyboardHelper()
	var kbd interface{}
	if kh != nil {
		kbd = kh.GetKeyboard()
	}
	txt += p.Resources.Text
	if kbd != nil {
		mr = &telegram.MessageRequest{ChatId: p.State.ChatId, Text: txt, ReplyMarkup: kbd, ParseMode: p.Resources.ParseMode}
	} else {
		mr = &telegram.MessageRequest{ChatId: p.State.ChatId, Text: txt, ParseMode: p.Resources.ParseMode}

	}
	return
}

func (p *ListdStateProvider) GetEditMR() (mer *telegram.EditMessageTextRequest) {
	mr := p.GetMR()
	mer = &telegram.EditMessageTextRequest{MessageId: p.State.MessageId, ChatId: mr.ChatId, Text: mr.Text, ParseMode: mr.ParseMode,
		ReplyMarkup: mr.ReplyMarkup}
	return
}

func (p *ListdStateProvider) InitReserves() {
	kh := telegram.NewDateKeyboardHelperRu()
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	if kh.Parse() != nil {
		return
	}
	dt := kh.Date
	filter := volley.Volley{}
	filter.StartTime = time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
	filter.EndTime = filter.StartTime.Add(time.Duration(time.Hour * 24))
	var err error
	if p.reserves, err = p.Repository.GetByFilter(filter, true, true); err != nil {
		log.WithFields(log.Fields{
			"package":  "bvbot",
			"function": "InitReserves",
			"struct":   "ListdStateProvider",
			"fliter":   filter,
			"error":    err,
		}).Error("can't get reserves by filter")
	}
}

func (p ListdStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.ChatId < 0 {
		return
	}
	if p.State.Action == "show" {
		return
	}
	p.InitReserves()
	rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetEditMR()})
	return
}

func (p ListdStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	if p.State.Action == "set" {
		ah := telegram.ActionsKeyboardHelper{}
		ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
		ah.Actions = []telegram.ActionButton{}

		ah.Columns = 1
		for i, res := range p.reserves {
			tgv := volley.NewTelegramViewRu(res)
			ab := telegram.ActionButton{
				Action: "show", Data: res.Base64Id(),
				Text: fmt.Sprintf("%v. %s", i+1, tgv.String())}
			ah.Actions = append(ah.Actions, ab)
		}
		return &ah
	}
	if p.State.Action == "listd" {
		kh := telegram.NewDateKeyboardHelperRu()
		kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
		return &kh
	}
	return nil
}

func (p ListdStateProvider) Proceed() (st telegram.State, err error) {
	return p.BaseStateProvider.Proceed()
}
