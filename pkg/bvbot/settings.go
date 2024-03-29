package bvbot

import (
	"strconv"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type SettingsStateProvider struct {
	BaseStateProvider
	Resources SettingsResources
}

func (p SettingsStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p SettingsStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	if p.reserve.Canceled {
		return &ah
	}
	ah.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
		if p.reserve.Person.TelegramId == p.Person.TelegramId || p.Person.CheckLocationRole(p.reserve.Location, "admin") {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "activity", Text: res.ActivityBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "courts", Text: res.CourtBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "level", Text: res.LevelBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "max", Text: res.MaxBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "price", Text: res.PriceBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "nettype", Text: res.NetTypeBtn})
		}
	}
	return &ah
}

type MaxPlayersStateProvider struct {
	BaseStateProvider
	Resources MaxPlayersResources
}

func (p MaxPlayersStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p MaxPlayersStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	cfg := p.GetLocationConfig().Courts
	kh := telegram.CountKeyboardHelper{
		Columns: res.Columns,
		Min:     cfg.MinPlayers,
		Max:     cfg.MaxPlayers * p.reserve.CourtCount,
		Step:    1,
	}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p MaxPlayersStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		if err := kh.Parse(); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "MaxPlayersStateProvider",
				"error":    err,
			}).Error("keyboard parse error")
		} else {
			p.reserve.MaxPlayers = kh.Count
			p.State.Updated = true
			p.State.Action = p.BackState.State
		}
	}
	return p.BaseStateProvider.Proceed()
}

type CourtsStateProvider struct {
	BaseStateProvider
	Resources CourtsResources
}

func (p CourtsStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p CourtsStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	cfg := p.GetLocationConfig().Courts
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: 1, Max: cfg.Max, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p CourtsStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		if err := kh.Parse(); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "CourtsStateProvider",
				"error":    err,
			}).Error("keyboard parse error")
		} else {
			p.reserve.CourtCount = kh.Count
			p.State.Updated = true
			p.State.Action = p.BackState.State
		}
	}
	return p.BaseStateProvider.Proceed()
}

type PriceStateProvider struct {
	BaseStateProvider
	Resources PriceResources
}

func (p PriceStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p PriceStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Price
	kh := telegram.CountKeyboardHelper{AlwaysZero: true, Columns: p.Resources.Columns,
		Min: cfg.Min, Max: cfg.Max, Step: cfg.Step}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(p.Resources.Message)
	return &kh
}

func (p PriceStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		if err := kh.Parse(); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "PriceStateProvider",
				"error":    err,
			}).Error("keyboard parse error")
		} else {
			p.reserve.Price = kh.Count
			p.State.Updated = true
			p.State.Action = p.BackState.State
		}
	}
	return p.BaseStateProvider.Proceed()
}

type ActivityStateProvider struct {
	BaseStateProvider
	Resources AcivityResources
}

func (p ActivityStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p ActivityStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources

	activities := []telegram.EnumItem{}
	for i := 0; i <= 30; i += 10 {
		activities = append(activities, telegram.EnumItem{Id: strconv.Itoa(i), Item: volley.Activity(i).String()})
	}
	kh := telegram.NewEnumKeyboardHelper(activities)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p ActivityStateProvider) Proceed() (telegram.State, error) {
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if p.State.Action == "set" {
		aid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ActivityStateProvider",
				"value":    kh.Value,
				"error":    err,
			}).Error("can't convert activity id")
		}
		p.reserve.Activity = volley.Activity(aid)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type LevelStateProvider struct {
	BaseStateProvider
	Resources LevelResources
}

func (p LevelStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p LevelStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources

	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i), Item: volley.PlayerLevel(i).String()})
	}
	kh := telegram.NewEnumKeyboardHelper(levels)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p LevelStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		aid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "LevelStateProvider",
				"value":    kh.Value,
				"error":    err,
			}).Error("can't convert level id")
		}
		p.reserve.MinLevel = aid
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type NetTypeStateProvider struct {
	BaseStateProvider
}

func (p NetTypeStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p NetTypeStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	types := []telegram.EnumItem{}
	for i := 0; i <= 20; i += 10 {
		types = append(types, telegram.EnumItem{Id: strconv.Itoa(i), Item: volley.NetType(i).String()})
	}
	kh := telegram.NewEnumKeyboardHelper(types)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p NetTypeStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		aid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "NetTypeStateProvider",
				"value":    kh.Value,
				"error":    err,
			}).Error("can't convert net type id")
		}
		p.reserve.NetType = volley.NetType(aid)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}
