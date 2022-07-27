package handlers

import (
	"fmt"
	"log"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/res"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type StateProvider interface {
	Init(telegram.User, telegram.Message, string) error
	ProceedCQ(telegram.User, telegram.CallbackQuery) error
	GetEditMR() telegram.EditMessageTextRequest
	GetMR() telegram.MessageRequest
	GetState() telegram.State
}

func NewReserveProvider(ps *services.PersonService, rr reserve.ReserveRepository, res *res.OrderResources) (rp ReserveProvider) {
	rp.ps = ps
	rp.rr = rr
	rp.res = res
	return
}

type ReserveProvider struct {
	res     *res.OrderResources
	ps      *services.PersonService
	rr      reserve.ReserveRepository
	kh      telegram.KeyboardHelper
	Reserve reserve.Reserve
	Person  person.Person
	User    telegram.User
	Message telegram.Message
}

func (s *ReserveProvider) Init(tu telegram.User, msg telegram.Message, data string) (err error) {
	s.User = tu
	s.Message = msg
	if s.Person, err = s.GetPerson(); err != nil {
		return
	}

	s.kh.SetData(data)
	if err = s.GetDataReserve(); err != nil {
		return
	}
	return
}

func (s *ReserveProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	s.User = tu
	s.Message = *cq.Message
	if s.Person, err = s.GetPerson(); err != nil {
		return
	}

	if err = s.kh.Parse(cq.Data); err != nil {
		return
	}

	if err = s.GetDataReserve(); err != nil {
		return
	}
	return
}

func (s *ReserveProvider) GetState() (st telegram.State) {
	st.State = s.kh.GetState()
	st.Data = s.Reserve.Id.String()
	st.ChatId = s.Message.Chat.Id
	st.MessageId = s.Message.MessageId
	return
}

func (s *ReserveProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	s.kh.SetData(s.Reserve.Id.String())
	return s.kh
}

func (s *ReserveProvider) GetPerson() (p person.Person, err error) {
	p, err = s.ps.GetByTelegramId(s.User.Id)
	if err != nil {
		log.Println(err.Error())
		_, ok := err.(person.ErrorPersonNotFound)
		if ok {
			p, _ = person.NewPerson(s.User.FirstName)
			s.Person.TelegramId = s.User.Id
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

func (s *ReserveProvider) GetDataReserve() (err error) {
	var id uuid.UUID
	id, err = uuid.Parse(s.kh.GetData())
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order getting reserve error: %s", err.Error()),
			AnswerMsg: "Parse reserve id error"}

	} else {
		s.Reserve, err = s.rr.Get(id)
		if err != nil {
			err = telegram.HelperError{
				Msg:       fmt.Sprintf("Getting reserve error: %s", err.Error()),
				AnswerMsg: "Getting reserve error"}
		}
	}
	return
}

func (s *ReserveProvider) GetEditMR() (mer telegram.EditMessageTextRequest) {
	mr := s.GetMR()
	return telegram.EditMessageTextRequest{ChatId: mr.ChatId, Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (s *ReserveProvider) GetMR() (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	s.kh.SetData(s.Reserve.Id.String())
	kbd.InlineKeyboard = append(kbd.InlineKeyboard, s.kh.GetKeyboard()...)
	kbdText = "\n*" + s.kh.GetText() + "* "

	rview := reserve.NewTelegramViewRu(s.Reserve)
	mtxt := fmt.Sprintf("%s\n%s", rview.GetText(), kbdText)
	if s.Message.Chat.Id < 0 {
		mtxt += s.res.MaxPlayer.GroupChatWarning
	}

	if len(kbd.InlineKeyboard) > 0 {
		return telegram.MessageRequest{ChatId: s.Message.Chat.Id, Text: mtxt, ParseMode: rview.ParseMode, ReplyMarkup: kbd}
	}
	return telegram.MessageRequest{ChatId: s.Message.Chat.Id, Text: mtxt, ParseMode: rview.ParseMode}
}

func (rm *ReserveProvider) UpdateReserve(res *reserve.Reserve) (err error) {

	if err = rm.rr.Update(*res); err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order update can't update reserve %s error: %s", res.Id, err.Error()),
			AnswerMsg: "Can't update reserve"}
		return
	}
	*res, err = rm.rr.Get(res.Id)
	return
}

type ReserveShowProvider struct {
	ReserveProvider
}

func NewReserveShowProvider(rp ReserveProvider) (rs *ReserveShowProvider) {
	rs = &ReserveShowProvider{ReserveProvider: rp}
	rs.kh = &telegram.ActionsKeyboardHelper{}
	return
}

func (s *ReserveShowProvider) Init(tu telegram.User, msg telegram.Message, data string) (err error) {
	if err = s.ReserveProvider.Init(tu, msg, data); err != nil {
		return
	}
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveShowProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	if err = s.ReserveProvider.ProceedCQ(tu, cq); err != nil {
		return
	}
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveShowProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	ah := s.kh.(*telegram.ActionsKeyboardHelper)
	ah.State = "ordershow"
	ah.Actions = []telegram.ActionButton{}
	ah.SetData(s.Reserve.Id.String())

	if s.Reserve.Canceled {
		return ah
	}
	ah.Columns = 2
	if s.Message.Chat.Id == s.Person.TelegramId {
		if s.Reserve.Person.TelegramId == s.Person.TelegramId || s.Person.CheckLocationRole(s.Reserve.Location, "admin") {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderdate", Text: s.res.DateTime.DateButton})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "ordertime", Text: s.res.DateTime.TimeButton})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "ordersets", Text: s.res.Set.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderdesc", Text: s.res.Description.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "ordersettings", Text: s.res.SettingsBtn})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderactions", Text: s.res.ActionsBtn})
		}
	}
	if s.Reserve.Ordered() {
		if s.Message.Chat.Id <= 0 || !s.Reserve.HasPlayerByTelegramId(s.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderjoin", Text: s.res.JoinPlayer.Button})
		}
		if s.Message.Chat.Id > 0 || s.Reserve.MaxPlayers-s.Reserve.PlayerCount(uuid.Nil) > 1 {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderjoinmult", Text: s.res.JoinPlayer.MultiButton})
		}
		if s.Message.Chat.Id > 0 && s.Reserve.HasPlayerByTelegramId(s.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderarrivetime", Text: s.res.JoinPlayer.ArriveButton})
		}
		if s.Message.Chat.Id <= 0 || s.Reserve.HasPlayerByTelegramId(s.Person.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Action: "orderleave", Text: s.res.JoinPlayer.LeaveButton})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordershow", Text: s.res.RefreshBtn})
	}

	return ah
}

type ReserveSettingsProvider struct {
	ReserveProvider
}

func NewReserveSettingsProvider(rp ReserveProvider) (rs *ReserveSettingsProvider) {
	rs = &ReserveSettingsProvider{ReserveProvider: rp}
	rs.kh = &telegram.ActionsKeyboardHelper{Columns: 2}
	return
}

func (s *ReserveSettingsProvider) Init(tu telegram.User, msg telegram.Message, data string) (err error) {
	if err = s.ReserveProvider.Init(tu, msg, data); err != nil {
		return
	}
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveSettingsProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	err = s.ReserveProvider.ProceedCQ(tu, cq)
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveSettingsProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	ah := s.kh.(*telegram.ActionsKeyboardHelper)
	ah.Actions = []telegram.ActionButton{}
	if s.Reserve.Person.TelegramId == s.Person.TelegramId || s.Person.CheckLocationRole(s.Reserve.Location, "admin") {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderactivity", Text: s.res.Activity.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordercourts", Text: s.res.Court.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderminlevel", Text: s.res.Level.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderplayers", Text: s.res.MaxPlayer.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderprice", Text: s.res.Price.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordershow", Text: s.res.BackBtn})
	}
	return ah
}

type ReserveActionsProvider struct {
	ReserveProvider
}

func NewReserveActionsProvider(rp ReserveProvider) (rs *ReserveActionsProvider) {
	rs = &ReserveActionsProvider{ReserveProvider: rp}
	rs.kh = &telegram.ActionsKeyboardHelper{Columns: 2}
	return
}

func (s *ReserveActionsProvider) Init(tu telegram.User, msg telegram.Message, data string) (err error) {
	if err = s.ReserveProvider.Init(tu, msg, data); err != nil {
		return
	}
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveActionsProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	err = s.ReserveProvider.ProceedCQ(tu, cq)
	s.kh = s.GetKeyboardHelper()
	return
}

func (s *ReserveActionsProvider) GetKeyboardHelper() telegram.KeyboardHelper {
	ah := s.kh.(*telegram.ActionsKeyboardHelper)
	ah.Actions = []telegram.ActionButton{}
	if s.Reserve.Person.TelegramId == s.Person.TelegramId || s.Person.CheckLocationRole(s.Reserve.Location, "admin") {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordercancel", Text: s.res.Cancel.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordercopy", Text: s.res.CopyBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderpub", Text: s.res.PublishBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "orderremovepl", Text: s.res.RemovePlayerBtn})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Action: "ordershow", Text: s.res.BackBtn})
	}
	return ah
}

type ReservePublishProvider struct {
	ReserveShowProvider
	Bot *telegram.Bot
}

func NewReservePublishProvider(rp ReserveShowProvider, tb *telegram.Bot) (rs *ReservePublishProvider) {
	rs = &ReservePublishProvider{ReserveShowProvider: rp, Bot: tb}
	return
}

func (s *ReservePublishProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	if err = s.ReserveShowProvider.ProceedCQ(tu, cq); err != nil {
		return
	}
	mr := s.GetMR()
	mr.ChatId = s.Reserve.Location.ChatId
	resp := s.Bot.SendMessage(&mr)
	s.Message = resp.Result
	return
}

type ReserveCopyProvider struct {
	ReserveShowProvider
}

func NewReserveCopyProvider(rp ReserveShowProvider) (rs *ReserveCopyProvider) {
	rs = &ReserveCopyProvider{ReserveShowProvider: rp}
	return
}

func (s *ReserveCopyProvider) ProceedCQ(tu telegram.User, cq telegram.CallbackQuery) (err error) {
	if err = s.ReserveShowProvider.ProceedCQ(tu, cq); err != nil {
		return
	}
	s.Reserve, err = s.rr.Add(s.Reserve.Copy())
	return
}

func (s *ReserveCopyProvider) GetEditMR() (mr telegram.EditMessageTextRequest) {
	mr = s.ReserveShowProvider.GetEditMR()
	mr.Text += s.res.CopyMessage
	return
}

func (s *ReserveCopyProvider) GetMR() (mr telegram.MessageRequest) {
	mr = s.ReserveShowProvider.GetMR()
	mr.Text += s.res.CopyMessage
	return
}

func (s *ReserveCopyProvider) GetState() (st telegram.State) {
	st = s.ReserveShowProvider.GetState()
	st.State = "ordercopy"
	return
}
