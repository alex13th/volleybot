package bvbot

import (
	"volleybot/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type ConfigCourtsStateProvider struct {
	ConfigStateProvider
}

func (p ConfigCourtsStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigCourtsStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 1
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgcourtmax", Text: res.Courts.MaxBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgcourtminpl", Text: res.Courts.MinPlayersBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "cfgcourtmaxpl", Text: res.Courts.MaxPlayersBtn})
	}
	return &ah
}

type ConfigCourtsMaxStateProvider struct {
	ConfigStateProvider
}

func (p ConfigCourtsMaxStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigCourtsMaxStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Courts
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: 1, Max: cfg.Max + 10, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigCourtsMaxStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Courts.Max = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigCourtsMaxStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}

type ConfigCourtsMinPlayersStateProvider struct {
	ConfigStateProvider
}

func (p ConfigCourtsMinPlayersStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigCourtsMinPlayersStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Courts
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: 1, Max: cfg.MaxPlayers, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigCourtsMinPlayersStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Courts.MinPlayers = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigCourtsMinPlayersStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}

type ConfigCourtsMaxPlayersStateProvider struct {
	ConfigStateProvider
}

func (p ConfigCourtsMaxPlayersStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.ConfigStateProvider.GetRequests()
}

func (p ConfigCourtsMaxPlayersStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	cfg := p.GetLocationConfig().Courts
	kh := telegram.CountKeyboardHelper{Columns: 3, Min: cfg.MinPlayers, Max: cfg.MaxPlayers + 10, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p ConfigCourtsMaxPlayersStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		cfg := p.GetLocationConfig()
		cfg.Courts.MaxPlayers = kh.Count
		p.State.Action = p.BackState.State
		if err := p.UpdateLocationConfig(cfg); err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "ConfigCourtsMaxPlayersStateProvider",
				"config":   cfg,
				"error":    err,
			}).Error("update location config error")
			return p.BackState, err
		}
	}
	return p.BaseStateProvider.Proceed()
}
