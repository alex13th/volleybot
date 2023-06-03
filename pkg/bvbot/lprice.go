package bvbot

import (
	"volleybot/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type ConfigPriceStateProvider struct {
	ConfigStateProvider
}

func (p ConfigPriceStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigPriceStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 1
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgpricemin", Text: res.Price.MinBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgpricemax", Text: res.Price.MaxBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgpricestep", Text: res.Price.StepBtn})
	}
	return &ah
}

type ConfigPriceMinStateProvider struct {
	ConfigStateProvider
}

func (p ConfigPriceMinStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigPriceMinStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Price
	min := cfg.Min - cfg.Step*5
	max := cfg.Max
	if min < 0 {
		min = 0
	}
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: min, Max: max, Step: cfg.Step}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigPriceMinStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Price.Min = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigPriceMinStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}

type ConfigPriceMaxStateProvider struct {
	ConfigStateProvider
}

func (p ConfigPriceMaxStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigPriceMaxStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Price
	min := cfg.Min
	max := cfg.Max + cfg.Step*5
	if min < 0 {
		min = 0
	}
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: min, Max: max, Step: cfg.Step}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigPriceMaxStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Price.Max = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigPriceMaxStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}

type ConfigPriceStepStateProvider struct {
	ConfigStateProvider
}

func (p ConfigPriceStepStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigPriceStepStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: 50, Max: 1000, Step: 50}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigPriceStepStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Price.Step = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigPriceStepStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}
