package handlers

import (
	"fmt"
	"log"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type ReserveProducer struct {
	res     *OrderResources
	ps      *services.PersonService
	rr      reserve.ReserveRepository
	kh      telegram.KeyboardHelper
	reserve reserve.Reserve
	person  person.Person
	User    telegram.User
	Chat    telegram.Chat
}

func (s *ReserveProducer) Init(tu telegram.User, tch telegram.Chat, data string) (err error) {
	s.User = tu
	s.Chat = tch
	if s.person, err = s.GetPerson(); err != nil {
		return
	}

	if err = s.kh.Parse(data); err != nil {
		return
	}

	if s.reserve, err = s.GetDataReserve(); err != nil {
		return
	}
	return
}

func (s *ReserveProducer) GetActions() telegram.KeyboardHelper {
	s.kh.SetData(s.reserve.Id.String())
	return s.kh
}

func (s *ReserveProducer) GetPerson() (p person.Person, err error) {
	p, err = s.ps.GetByTelegramId(s.User.Id)
	if err != nil {
		log.Println(err.Error())
		_, ok := err.(person.ErrorPersonNotFound)
		if ok {
			p, _ = person.NewPerson(s.User.FirstName)
			s.person.TelegramId = s.User.Id
			p.Lastname = s.User.LastName
			p, err = s.ps.Add(p)
		}
	}
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("getting person error: %s", err.Error()),
			AnswerMsg: "Can't get person"}

	}

	return
}

func (s *ReserveProducer) GetDataReserve() (r reserve.Reserve, err error) {
	var id uuid.UUID
	id, err = uuid.Parse(s.kh.GetData())
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order getting reserve error: %s", err.Error()),
			AnswerMsg: "Parse reserve id error"}

	} else {
		r, err = s.rr.Get(id)
		if err != nil {
			err = telegram.HelperError{
				Msg:       fmt.Sprintf("Getting reserve error: %s", err.Error()),
				AnswerMsg: "Getting reserve error"}
		}
	}
	return
}

func (s *ReserveProducer) GetEditMR() (mer telegram.EditMessageTextRequest) {
	mr := s.GetMR()
	return telegram.EditMessageTextRequest{ChatId: mr.ChatId, Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (s *ReserveProducer) GetMR() (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	s.kh.SetData(s.reserve.Id.String())
	kbd.InlineKeyboard = append(kbd.InlineKeyboard, s.kh.GetKeyboard()...)
	kbdText = "\n*" + s.kh.GetText() + "* "

	rview := reserve.NewTelegramViewRu(s.reserve)
	mtxt := fmt.Sprintf("%s\n%s", rview.GetText(), kbdText)
	if s.Chat.Id < 0 {
		mtxt += s.res.MaxPlayer.GroupChatWarning
	}

	if len(kbd.InlineKeyboard) > 0 {
		return telegram.MessageRequest{ChatId: s.Chat.Id, Text: mtxt, ParseMode: rview.ParseMode, ReplyMarkup: kbd}
	}
	return telegram.MessageRequest{ChatId: s.Chat.Id, Text: mtxt, ParseMode: rview.ParseMode}
}

type ReserveShowProucer struct {
	ReserveProducer
}

func NewReserveShowProducer(ps *services.PersonService, rr reserve.ReserveRepository, res *OrderResources) (rs ReserveShowProucer) {
	rs.res = res
	rs.ps = ps
	rs.rr = rr
	rs.kh = &telegram.ActionsKeyboardHelper{}
	return
}

func (s *ReserveShowProucer) Init(tu telegram.User, tch telegram.Chat, data string) (err error) {
	err = s.ReserveProducer.Init(tu, tch, data)
	s.kh = s.GetActions()
	return
}

func (s *ReserveShowProucer) GetActions() telegram.KeyboardHelper {
	ah := s.ReserveProducer.GetActions().(*telegram.ActionsKeyboardHelper)
	ah.SetData(s.reserve.Id.String())

	if s.reserve.Canceled {
		return ah
	}
	ah.Columns = 2
	if s.Chat.Id == s.person.TelegramId {
		if s.reserve.Person.TelegramId == s.person.TelegramId || s.person.CheckLocationRole(s.reserve.Location, "admin") {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderdate", Text: s.res.DateTime.DateButton})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "ordertime", Text: s.res.DateTime.TimeButton})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "ordersets", Text: s.res.Set.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderdesc", Text: s.res.Description.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "ordersettings", Text: s.res.SettingsBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderactions", Text: s.res.ActionsBtn})
		}
	}
	if s.reserve.Ordered() {
		if s.Chat.Id <= 0 || !s.reserve.HasPlayerByTelegramId(s.person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderjoin", Text: s.res.JoinPlayer.Button})
		}
		if s.Chat.Id > 0 || s.reserve.MaxPlayers-s.reserve.PlayerCount(uuid.Nil) > 1 {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderjoinmult", Text: s.res.JoinPlayer.MultiButton})
		}
		if s.Chat.Id > 0 && s.reserve.HasPlayerByTelegramId(s.person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderarrivetime", Text: s.res.JoinPlayer.ArriveButton})
		}
		if s.Chat.Id <= 0 || s.reserve.HasPlayerByTelegramId(s.person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderleave", Text: s.res.JoinPlayer.LeaveButton})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordershow", Text: s.res.RefreshBtn})
	}

	return ah
}

type ReserveSettingsProducer struct {
	ReserveProducer
}

func NewReserveSettingsProducer(ps *services.PersonService, rr reserve.ReserveRepository, res *OrderResources) (rs ReserveSettingsProducer) {
	rs.res = res
	rs.ps = ps
	rs.rr = rr
	rs.kh = &telegram.ActionsKeyboardHelper{Columns: 2}
	return
}

func (s *ReserveSettingsProducer) Init(tu telegram.User, tch telegram.Chat, data string) (err error) {
	err = s.ReserveProducer.Init(tu, tch, data)
	s.kh = s.GetActions()
	return
}

func (s *ReserveSettingsProducer) GetActions() telegram.KeyboardHelper {
	ah := s.ReserveProducer.GetActions().(*telegram.ActionsKeyboardHelper)
	if s.reserve.Person.TelegramId == s.person.TelegramId || s.person.CheckLocationRole(s.reserve.Location, "admin") {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderactivity", Text: s.res.Activity.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordercourts", Text: s.res.Court.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderminlevel", Text: s.res.Level.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderplayers", Text: s.res.MaxPlayer.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderprice", Text: s.res.Price.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordershow", Text: s.res.BackBtn})
	}
	return ah
}
