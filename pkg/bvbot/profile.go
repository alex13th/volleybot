package bvbot

import (
	"fmt"
	"strconv"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/volley"

	"volleybot/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

type PlayerStateProvider struct {
	BaseStateProvider
	Resources ProfileResources
	Player    volley.Player
}

func (p PlayerStateProvider) Proceed() (st telegram.State, err error) {
	if p.State.Updated {
		p.Repository.UpdatePlayer(p.Player)
		p.State.Updated = false
	}
	return p.BaseStateProvider.Proceed()
}

func (p PlayerStateProvider) GetMR() *telegram.MessageRequest {
	pview := volley.NewPlayerTelegramView(p.Player)
	psetview := person.NewTelegramSettingsViewRu(p.Player.Person)
	txt := fmt.Sprintf("%s\n\n%s\n%s", pview.GetText(), psetview.GetText(), p.Text)

	kbd := p.kh.GetKeyboard()

	return p.CreateMR(p.State.ChatId, txt, p.Resources.ParseMode, kbd)
}

func (p PlayerStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	var err error
	if p.Player, err = p.Repository.GetPlayer(p.Person); err != nil {
		log.WithFields(log.Fields{
			"package":  "bvbot",
			"function": "GetRequests",
			"struct":   "PlayerStateProvider",
			"person":   p.Person,
			"error":    err,
		}).Error("can't get player information")
	}
	var sreq telegram.StateRequest
	sreq.State = p.State
	sreq.Request = p.GetEditMR(p.GetMR())
	return append(reqlist, sreq)
}

type ProfileStateProvider struct {
	PlayerStateProvider
}

func (p ProfileStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	p.kh = p.GetKeyboardHelper()
	return p.PlayerStateProvider.GetRequests()
}

func (p ProfileStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "plevel", Text: res.LevelBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "sex", Text: res.SexBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "notifies", Text: res.NotifiesBtn})
	}
	return &ah
}

type PLevelStateProvider struct {
	PlayerStateProvider
	Resources LevelResources
}

func (p PLevelStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.PlayerStateProvider.GetRequests()
}

func (p PLevelStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources

	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i), Item: volley.PlayerLevel(i).String()})
	}
	kh := telegram.NewEnumKeyboardHelper(levels)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p PLevelStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		aid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "PLevelStateProvider",
				"value":    kh.Value,
				"error":    err,
			}).Error("can't parse level value")
		}
		p.Player, err = p.Repository.GetPlayer(p.Person)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "PLevelStateProvider",
				"person":   p.Person,
				"error":    err,
			}).Error("can't get player information")
		}
		p.Player.Level = volley.PlayerLevel(aid)
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.PlayerStateProvider.Proceed()
}

type SexStateProvider struct {
	PlayerStateProvider
}

func (p SexStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.PlayerStateProvider.GetRequests()
}

func (p SexStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	sexs := []telegram.EnumItem{
		{Id: "1", Item: fmt.Sprintf("%s %s", person.Sex(1).Emoji(), person.Sex(1))},
		{Id: "2", Item: fmt.Sprintf("%s %s", person.Sex(2).Emoji(), person.Sex(2))},
	}

	kh := telegram.NewEnumKeyboardHelper(sexs)
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

func (p SexStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		aid, err := strconv.Atoi(kh.Value)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "SexStateProvider",
				"value":    kh.Value,
				"error":    err,
			}).Error("can't parse sex value")
		}
		p.Player, err = p.Repository.GetPlayer(p.Person)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "SexStateProvider",
				"person":   p.Person,
				"error":    err,
			}).Error("can't get player information")
		}
		p.Player.Sex = person.Sex(aid)
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.PlayerStateProvider.Proceed()
}

type NotifiesStateProvider struct {
	PlayerStateProvider
}

func (p NotifiesStateProvider) GetRequests() (reqlist []telegram.StateRequest) {
	p.kh = p.GetKeyboardHelper()
	return p.PlayerStateProvider.GetRequests()
}

func (p NotifiesStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	ah.Columns = 2
	if p.State.ChatId == p.Person.TelegramId {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "pcancel", Text: res.CancelNotifyBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "pnotify", Text: res.NotifyBtn})
	}
	return &ah
}

type ParamStateProvider struct {
	PlayerStateProvider
}

func (p ParamStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.PlayerStateProvider.GetRequests()
}

func (p ParamStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	bparams := []telegram.EnumItem{
		{Id: "undef", Item: person.ParamValText["undef"]},
		{Id: "on", Item: person.ParamValText["on"]},
		{Id: "off", Item: person.ParamValText["off"]},
	}

	kh := telegram.NewEnumKeyboardHelper(bparams)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	return &kh
}

type NotifyCancelStateProvider struct {
	ParamStateProvider
}

func (p NotifyCancelStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		var err error
		p.Player, err = p.Repository.GetPlayer(p.Person)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "NotifyCancelStateProvider",
				"person":   p.Person,
				"error":    err,
			}).Error("can't get player information")
		}
		p.Player.Settings["notify_cancel"] = kh.Value
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.PlayerStateProvider.Proceed()
}

type NotifyChangeStateProvider struct {
	ParamStateProvider
}

func (p NotifyChangeStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if st.Action == "set" {
		var err error
		p.Player, err = p.Repository.GetPlayer(p.Person)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "bvbot",
				"function": "Proceed",
				"struct":   "NotifyChangeStateProvider",
				"person":   p.Person,
				"error":    err,
			}).Error("can't get player information")
		}
		p.Player.Settings["notify"] = kh.Value
		p.State.Action = p.BackState.State
		p.State.Updated = true
	}
	return p.PlayerStateProvider.Proceed()
}
