package bvbot

import (
	"strconv"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"
)

type SettingsResources struct {
	ActivityBtn string
	BackBtn     string
	CourtBtn    string
	LevelBtn    string
	MaxBtn      string
	NetTypeBtn  string
	PriceBtn    string
}

func NewSettingsResourcesRu() (r SettingsResources) {
	r.ActivityBtn = "Вид активности"
	r.BackBtn = "Назад"
	r.CourtBtn = "🏐 Площадки"
	r.LevelBtn = "💪 Уровень"
	r.MaxBtn = "👫 Мест"
	r.NetTypeBtn = "📏 Вид сетки"
	r.PriceBtn = "💰 Стоимость"
	return
}

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

type MaxPlayersResources struct {
	BackBtn          string `json:"back_btn"`
	Columns          int    `json:"columns"`
	Courts           CourtsResources
	GroupChatWarning string `json:"group_chat_warning"`
	Message          string `json:"message"`
}

func NewMaxPlayersResourcesRu() (r MaxPlayersResources) {
	r.BackBtn = "Назад"
	r.Columns = 4
	r.Courts = NewCourtsResourcesRu()
	r.GroupChatWarning = "⚠️*Внимание* - здесь функция добавления игроков ограничена числом игроков записи. " +
		"В чате с ботом можно добавить больше игроков в резерв!"
	return
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
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: 4, Max: res.Courts.MaxPlayers * p.reserve.CourtCount, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p MaxPlayersStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		p.reserve.MaxPlayers = kh.Count
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type CourtsResources struct {
	Columns    int    `json:"columns"`
	Max        int    `json:"max"`
	MaxPlayers int    `json:"max_players"`
	Message    string `json:"message"`
}

func NewCourtsResourcesRu() CourtsResources {
	return CourtsResources{Columns: 4, Max: 4, Message: "❓Сколько нужно кортов❓", MaxPlayers: 12}
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
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: 1, Max: res.Max, Step: 1}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p CourtsStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		p.reserve.CourtCount = kh.Count
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type PriceResources struct {
	Columns int
	Min     int
	Max     int
	Message string
	Step    int
}

func NewPriceResourcesRu() PriceResources {
	return PriceResources{Columns: 4, Min: 0, Max: 2000, Message: "❓Почем будет поиграть❓", Step: 100}
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
	res := p.Resources
	kh := telegram.CountKeyboardHelper{Columns: res.Columns, Min: res.Min, Max: res.Max, Step: res.Step}
	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p PriceStateProvider) Proceed() (telegram.State, error) {
	st := p.State
	kh := p.GetKeyboardHelper().(*telegram.CountKeyboardHelper)
	if st.Action == "set" {
		kh.Parse()
		p.reserve.Price = kh.Count
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type AcivityResources struct {
	Columns int    `json:"columns"`
	Message string `json:"message"`
}

func NewAcivityResourcesRu() AcivityResources {
	return AcivityResources{Columns: 1, Message: "❓Какой будет вид активности❓"}
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
		activities = append(activities, telegram.EnumItem{Id: strconv.Itoa(i), Item: reserve.Activity(i).String()})
	}
	kh := telegram.NewEnumKeyboardHelper(activities)

	kh.BaseKeyboardHelper = p.GetBaseKeyboardHelper(res.Message)
	return &kh
}

func (p ActivityStateProvider) Proceed() (telegram.State, error) {
	kh := p.GetKeyboardHelper().(*telegram.EnumKeyboardHelper)
	if p.State.Action == "set" {
		aid, _ := strconv.Atoi(kh.Value)
		p.reserve.Activity = reserve.Activity(aid)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}

type LevelResources struct {
	Columns int    `json:"columns"`
	Message string `json:"message"`
}

func NewLevelResourcesRu() LevelResources {
	return LevelResources{Columns: 3, Message: "❓Какой минимальный уровень игроков❓"}
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
		aid, _ := strconv.Atoi(kh.Value)
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
		aid, _ := strconv.Atoi(kh.Value)
		p.reserve.NetType = volley.NetType(aid)
		p.State.Updated = true
		p.State.Action = p.BackState.State
	}
	return p.BaseStateProvider.Proceed()
}
