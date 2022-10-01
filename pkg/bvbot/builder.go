package bvbot

import (
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"
)

type BaseStateProvider struct {
	reserve          volley.Volley
	kh               telegram.KeyboardHelper
	name             string
	BackState        telegram.State
	Message          telegram.Message
	Person           person.Person
	Repository       volley.Repository
	ConfigRepository location.LocationConfigRepository
	Location         location.Location
	State            telegram.State
	Text             string
}

func NewBaseStateProvider(state telegram.State, msg telegram.Message, p person.Person, loc location.Location,
	rep volley.Repository, cfgrep location.LocationConfigRepository, text string) (sp BaseStateProvider, err error) {
	sp = BaseStateProvider{State: state, Message: msg, Person: p, Location: loc, Repository: rep, ConfigRepository: cfgrep, Text: text}
	sp.name = "beach_volley"
	if rep != nil && state.Data != "" {
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
	if p.State.ChatId > 0 {
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
	if p.State.ChatId < 0 {
		req := telegram.DeleteMessageRequest{ChatId: p.State.ChatId, MessageId: p.State.MessageId}
		rlist = append(rlist, telegram.StateRequest{Request: &req})
		rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetMR()})

	} else {
		rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetEditMR(p.GetMR())})
	}
	return
}

func (p *BaseStateProvider) GetEditMR(mr *telegram.MessageRequest) (mer *telegram.EditMessageTextRequest) {
	mer = &telegram.EditMessageTextRequest{MessageId: p.State.MessageId, ChatId: mr.ChatId, Text: mr.Text, ParseMode: mr.ParseMode,
		ReplyMarkup: mr.ReplyMarkup}
	return
}

func (p *BaseStateProvider) GetMR() (mr *telegram.MessageRequest) {
	cid := p.State.ChatId
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
	return p.CreateMR(cid, mtxt, rview.ParseMode, kbd)
}

func (p BaseStateProvider) CreateMR(cid int, txt string, pmode string, kbd telegram.InlineKeyboardMarkup) *telegram.MessageRequest {
	if len(kbd.InlineKeyboard) > 0 {
		return &telegram.MessageRequest{ChatId: cid, Text: txt, ParseMode: pmode,
			ReplyMarkup: kbd, DisableNotification: cid < 0}
	}
	return &telegram.MessageRequest{ChatId: cid, Text: txt, ParseMode: pmode,
		DisableNotification: cid < 0}
}

func (p *BaseStateProvider) NotifyPlayers(action string) (reqlist []telegram.StateRequest) {
	for _, mb := range p.reserve.Members {
		if mb.Person.TelegramId != p.Person.TelegramId {
			if param, ok := mb.Settings[action]; ok && param == "on" {
				mr := p.GetMR()
				mr.ChatId = mb.TelegramId
				reqlist = append(reqlist, telegram.StateRequest{Request: mr})
			}
		}
	}
	return
}

func (p BaseStateProvider) GetLocationConfig() (conf Config) {
	err := p.ConfigRepository.Get(p.Location, p.name, &conf)

	if err != nil {
		if (conf == Config{}) {
			conf = NewConfig()
			p.ConfigRepository.Add(p.Location, p.name, conf)
		} else {
			p.ConfigRepository.Update(p.Location, p.name, conf)
		}
	}

	return
}

func (p BaseStateProvider) UpdateLocationConfig(conf Config) error {
	return p.ConfigRepository.Update(p.Location, p.name, &conf)
}

type BvStateBuilder struct {
	BaseStateProvider
	Resources Resources
}

func NewBvStateBuilder(loc location.Location, msg telegram.Message, p person.Person, rep volley.Repository, res Resources, cfgrep location.LocationConfigRepository, st telegram.State) (bld BvStateBuilder, err error) {
	bp, err := NewBaseStateProvider(st, msg, p, loc, rep, cfgrep, "")
	bld = BvStateBuilder{BaseStateProvider: bp, Resources: res}
	return
}

func (bld BvStateBuilder) GetStateProvider(st telegram.State) (sp telegram.StateProvider, err error) {
	bp := bld.BaseStateProvider
	bp.State = st
	bp.BackState = st
	switch bp.State.State {
	case "main":
		bp.BackState = telegram.State{}
		sp = MainStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Main}
	case "config":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		sp = ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
	case "listd":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		sp = ListdStateProvider{BaseStateProvider: bp, Resources: bld.Resources.List}
	case "profile":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = ProfileStateProvider{PlayerStateProvider: pp}
	case "show":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		sp = ShowStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Show}
	case "actions":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = ActionsStateProvider{BaseStateProvider: bp,
			Resources: bld.Resources.Actions, ShowResources: bld.Resources.Show}
	case "date":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = DateStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Show.DateTime}
	case "desc":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = &DescStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Description}
	case "time":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = TimeStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Show.DateTime}
	case "sets":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = SetsStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Sets}
	case "joinm":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = JoinPlayersStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Join}
	case "jtime":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = JoinTimeStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Join}
	case "settings":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = SettingsStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Settings}
	case "courts":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = CourtsStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Courts}
	case "max":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = MaxPlayersStateProvider{BaseStateProvider: bp, Resources: bld.Resources.MaxPlayer}
	case "price":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = PriceStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Price}
	case "level":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = LevelStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Level}
	case "activity":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = ActivityStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Activity}
	case "nettype":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = NetTypeStateProvider{BaseStateProvider: bp}
	case "cancel":
		bp.BackState.State = "actions"
		bp.BackState.Action = bp.BackState.State
		sp = CancelStateProvider{BaseStateProvider: bp,
			Resources: bld.Resources.Cancel, ShowResources: bld.Resources.Show}
	case "rmpl":
		bp.BackState.State = "actions"
		bp.BackState.Action = bp.BackState.State
		sp = RemovePlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.RemovePlayer}
	case "plevel":
		bp.BackState.State = "profile"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = PLevelStateProvider{PlayerStateProvider: pp, Resources: bld.Resources.Level}
	case "sex":
		bp.BackState.State = "profile"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = SexStateProvider{PlayerStateProvider: pp}
	case "notifies":
		bp.BackState.State = "profile"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = NotifiesStateProvider{PlayerStateProvider: pp}
	case "pcancel":
		bp.BackState.State = "notifies"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = NotifyCancelStateProvider{ParamStateProvider{PlayerStateProvider: pp}}
	case "pnotify":
		bp.BackState.State = "notifies"
		bp.BackState.Action = bp.BackState.State
		pp := PlayerStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Profile}
		sp = NotifyChangeStateProvider{ParamStateProvider{PlayerStateProvider: pp}}
	case "cfgcourts":
		bp.BackState.State = "config"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigCourtsStateProvider{ConfigStateProvider: cfgp}
	case "cfgcourtmax":
		bp.BackState.State = "cfgcourts"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigCourtsMaxStateProvider{ConfigStateProvider: cfgp}
	case "cfgcourtminpl":
		bp.BackState.State = "cfgcourts"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigCourtsMinPlayersStateProvider{ConfigStateProvider: cfgp}
	case "cfgcourtmaxpl":
		bp.BackState.State = "cfgcourts"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigCourtsMaxPlayersStateProvider{ConfigStateProvider: cfgp}
	case "cfgprice":
		bp.BackState.State = "config"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigPriceStateProvider{ConfigStateProvider: cfgp}
	case "cfgpricemin":
		bp.BackState.State = "cfgprice"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigPriceMinStateProvider{ConfigStateProvider: cfgp}
	case "cfgpricemax":
		bp.BackState.State = "cfgprice"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigPriceMaxStateProvider{ConfigStateProvider: cfgp}
	case "cfgpricestep":
		bp.BackState.State = "cfgprice"
		bp.BackState.Action = bp.BackState.State
		cfgp := ConfigStateProvider{BaseStateProvider: bp, Resources: bld.Resources.Config}
		sp = ConfigPriceStepStateProvider{ConfigStateProvider: cfgp}
	}

	return
}
