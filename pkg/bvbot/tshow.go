package bvbot

import (
	"volleybot/pkg/domain/person"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type ShowResources struct {
	DateTime       telegram.DateTimeResources
	ActionsBtn     string
	DescriptionBtn string
	JoinBtn        string
	JoinLeaveBtn   string
	JoinMultiBtn   string
	JoinTimeBtn    string
	RefreshBtn     string
	SetsBtn        string
	SettingsBtn    string
}

func NewShowResourcesRu() (r ShowResources) {
	r.DateTime = telegram.NewDateTimeResourcesRu()
	r.ActionsBtn = "Действия"
	r.DescriptionBtn = "Описание"
	r.JoinBtn = "😀 Буду"
	r.JoinLeaveBtn = "😞 Не смогу"
	r.JoinMultiBtn = "🤩 Буду не один"
	r.JoinTimeBtn = "🏃‍♂️ Опоздаю"
	r.RefreshBtn = "Обновить"
	r.SetsBtn = "⏱ Кол-во часов"
	r.SettingsBtn = "Настройки"
	return
}

type ShowStateProvider struct {
	BaseStateProvider
	Resources ShowResources
}

func (p ShowStateProvider) GetRequests() []telegram.StateRequest {
	if p.State.Action != "show" {
		return nil
	}
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p ShowStateProvider) GetKeyboardHelper() (kh telegram.KeyboardHelper) {
	res := p.Resources
	msg := p.Message
	ah := telegram.ActionsKeyboardHelper{}
	ah.BaseKeyboardHelper = p.GetBaseKeyboardHelper("")
	ah.Actions = []telegram.ActionButton{}

	if p.reserve.Canceled {
		return &ah
	}
	ah.Columns = 2
	if msg.Chat.Id == p.Person.TelegramId {
		if p.reserve.Person.TelegramId == p.Person.TelegramId || p.Person.CheckLocationRole(p.reserve.Location, "admin") {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "date", Text: res.DateTime.DateBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "time", Text: res.DateTime.TimeBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "sets", Text: res.SetsBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "desc", Text: res.DescriptionBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "settings", Text: res.SettingsBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "actions", Text: res.ActionsBtn})
		}
	}
	if p.reserve.Ordered() {
		if msg.Chat.Id <= 0 || !p.reserve.HasPlayerByTelegramId(p.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "join", Text: res.JoinBtn})
		}
		if msg.Chat.Id > 0 || p.reserve.MaxPlayers-p.reserve.PlayerCount(uuid.Nil) > 1 {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "joinm", Text: res.JoinMultiBtn})
		}
		if msg.Chat.Id > 0 && p.reserve.HasPlayerByTelegramId(p.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "jtime", Text: res.JoinTimeBtn})
		}
		if msg.Chat.Id <= 0 || p.reserve.HasPlayerByTelegramId(p.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "leave", Text: res.JoinLeaveBtn})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "refresh", Text: res.RefreshBtn})
	}

	return &ah

}

func (p ShowStateProvider) Proceed() (telegram.State, error) {
	if p.State.Action == "refresh" {
		p.State.Action = "show"
	}
	if p.State.Action == "join" {
		rpl := p.reserve.GetPlayer(p.Person.Id)
		if rpl.Id == uuid.Nil {
			rpl = person.Player{Person: p.Person}
		}
		rpl.Count = 1
		p.reserve.JoinPlayer(rpl)
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
		p.State.Action = "show"
		p.State.Updated = true
	}
	return p.BaseStateProvider.Proceed()
}

type DateStateProvider struct {
	BaseStateProvider
	Resources telegram.DateTimeResources
}

func (p DateStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p DateStateProvider) Proceed() (telegram.State, error) {
	kh := p.GetKeyboardHelper().(*telegram.DateKeyboardHelper)
	if p.State.Action == "set" {
		kh.Parse()
		p.reserve.SetStartDate(kh.Date)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

func (p DateStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.NewDateKeyboardHelperRu()
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.DateMsg)
	return &kh
}

type TimeStateProvider struct {
	BaseStateProvider
	Resources telegram.DateTimeResources
}

func (p TimeStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p TimeStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.NewTimeKeyboardHelperRu()
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.TimeMsg)
	return &kh
}

func (p TimeStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.TimeKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		p.reserve.SetStartTime(kh.Time)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

func NewSetsResourcesRu() (r SetsResources) {
	r.BackBtn = "Назад"
	r.Columns = 4
	r.Max = 14
	r.Message = "❓Количество часов❓"
	return
}

type SetsResources struct {
	BackBtn string
	Columns int
	Max     int
	Message string
}

type SetsStateProvider struct {
	BaseStateProvider
	Resources SetsResources
}

func (p SetsStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p SetsStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: 1, Max: res.Max, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p SetsStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		p.reserve.SetDurationHours(kh.Count)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type JoinResources struct {
	BackBtn  string
	Columns  int
	Courts   CourtsResources
	DateTime telegram.DateTimeResources
	Message  string
}

func NewJoinPlayersResourcesRu() (r JoinResources) {
	r.BackBtn = "Назад"
	r.Columns = 4
	r.Courts = NewCourtsResourcesRu()
	r.Message = "❓Сколько игроков записать❓"
	r.DateTime = telegram.NewDateTimeResourcesRu()
	return
}

type JoinPlayersStateProvider struct {
	BaseStateProvider
	Resources JoinResources
}

func (p JoinPlayersStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action == "joinm" {
		p.kh = p.GetKeyboardHelper()
		rlist = append(rlist, telegram.StateRequest{State: p.State, Request: p.GetEditMR()})
	}
	return
}

func (p JoinPlayersStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	msg := p.Message
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: 1, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)

	kh.Min = 1
	if msg.Chat.Id > 0 {
		kh.Max = p.reserve.MaxPlayers - p.reserve.PlayerCount(p.Person.Id)
		if !p.reserve.HasPlayerByTelegramId(p.Person.TelegramId) || p.reserve.PlayerInReserve(p.Person.Id) {
			kh.Max = p.reserve.MaxPlayers
		} else if kh.Max <= p.reserve.GetPlayer(p.Person.Id).Count {
			kh.Max = p.reserve.GetPlayer(p.Person.Id).Count
		}
	} else {
		kh.Max = p.reserve.MaxPlayers - p.reserve.PlayerCount(uuid.Nil)
	}
	return &kh
}

func (p JoinPlayersStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		rpl := p.reserve.GetPlayer(p.Person.Id)
		if rpl.Id == uuid.Nil {
			rpl = person.Player{Person: p.Person}
		}
		rpl.Count = kh.Count
		p.reserve.JoinPlayer(rpl)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type JoinTimeStateProvider struct {
	BaseStateProvider
	Resources JoinResources
}

func (p JoinTimeStateProvider) GetRequests() []telegram.StateRequest {
	p.kh = p.GetKeyboardHelper()
	return p.BaseStateProvider.GetRequests()
}

func (p JoinTimeStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	res := p.Resources
	kh := telegram.NewTimeKeyboardHelperRu()
	kh.Step = 15
	kh.StartHour = p.reserve.StartTime.Hour()
	kh.EndHour = p.reserve.EndTime.Hour() - 1
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.DateTime.TimeMsg)
	return &kh
}

func (p JoinTimeStateProvider) Proceed() (telegram.State, error) {
	if p.State.Action == "set" {
		kh := p.GetKeyboardHelper().(*telegram.TimeKeyboardHelper)
		kh.Parse()
		pl := p.reserve.GetPlayerByTelegramId(p.Person.TelegramId)
		pl.ArriveTime = kh.Time
		p.reserve.JoinPlayer(pl)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type DescResources struct {
	BackBtn     string
	Message     string
	DoneMessage string
}

func NewDescResourcesRu() (r DescResources) {
	r.BackBtn = "Назад"
	r.Message = "Отлично. Отправь в чат описание активности."
	r.DoneMessage = "Успешно! Описание обновлено."
	return
}

type DescStateProvider struct {
	BaseStateProvider
	Resources DescResources
}

func (p DescStateProvider) GetRequests() (rlist []telegram.StateRequest) {
	if p.State.Action == "done" {
		rlist = append(rlist, telegram.StateRequest{Clear: true, State: p.State})
		req := telegram.MessageRequest{ChatId: p.Message.Chat.Id, Text: p.Resources.DoneMessage}
		return append(rlist, telegram.StateRequest{Request: &req})
	}
	if p.State.Action == "desc" {
		req := telegram.MessageRequest{ChatId: p.Message.Chat.Id, Text: p.Resources.Message}
		p.State.MessageId = -1
		return append(rlist, telegram.StateRequest{State: p.State, Request: &req})
	}
	return
}

func (p DescStateProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	return nil
}

func (p *DescStateProvider) Proceed() (telegram.State, error) {
	if p.State.Action == "desc" {
		p.reserve.Description = p.Message.Text
		err := p.Repository.Update(p.reserve)
		p.State.Action = "done"
		p.BackState.Updated = true
		return p.BackState, err
	}
	return p.State, nil
}
