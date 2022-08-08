package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/res"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/goodsign/monday"
	"github.com/google/uuid"
)

func NewOrderHandler(tb *telegram.Bot, os *services.OrderService, sr telegram.StateRepository, rl res.OrderResourceLoader) (oh OrderBotHandler) {
	oh = OrderBotHandler{OrderService: os}
	oh.Resources = rl.GetResource()
	oh.CommonHandler = CommonHandler{
		Bot:             tb,
		StateRepository: sr,
		PersonService:   os.PersonService,
		Resources:       oh.Resources,
	}
	oh.PlayerHandler = &PlayerHandler{
		PersonService: os.PersonService,
		CommonHandler: oh.CommonHandler,
	}
	oh.ReserveHandler = &ReserveHandler{
		CommonHandler: oh.CommonHandler,
		PlayerHandler: oh.PlayerHandler,
		Reserves:      os.Reserves,
	}

	oh.DateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderdate")
	oh.DateHelper.Days = oh.Resources.DateTime.DayCount
	oh.DateHelper.Columns = 3

	oh.ListDateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderlistdate")
	oh.ListDateHelper.Days = oh.Resources.DateTime.DayCount
	oh.ListDateHelper.Columns = 3

	oh.TimeHelper = telegram.NewTimeKeyboardHelper(oh.Resources.DateTime.TimeMessage, "ordertime")

	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i), Item: person.PlayerLevel(i).String()})
	}
	oh.MinLevelHelper = telegram.NewEnumKeyboardHelper(oh.Resources.Level.Message, "orderminlevel", levels)
	activities := []telegram.EnumItem{}
	for i := 0; i <= 30; i += 10 {
		activities = append(activities, telegram.EnumItem{Id: strconv.Itoa(i), Item: reserve.Activity(i).String()})
	}
	oh.ActivityHelper = telegram.NewEnumKeyboardHelper(oh.Resources.Activity.Message, "orderactivity", activities)

	oh.CourtsHelper = telegram.NewCountKeyboardHelper(oh.Resources.Court.Message, "ordercourts", 1, oh.Resources.Court.Max)
	oh.SetsHelper = telegram.NewCountKeyboardHelper(oh.Resources.Set.Message, "ordersets", 1, oh.Resources.Court.Max)
	oh.PlayerCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.MaxPlayer.Message, "orderplayers", oh.Resources.MaxPlayer.Min, oh.Resources.MaxPlayer.Max)
	oh.JoinCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.JoinPlayer.Message, "orderjoinmult", oh.Resources.MaxPlayer.Min, oh.Resources.MaxPlayer.Max)
	oh.PriceCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.Price.Message, "orderprice", oh.Resources.Price.Min, oh.Resources.Price.Max)
	oh.PriceCountHelper.Step = oh.Resources.Price.Step
	return
}

type OrderBotHandler struct {
	CommonHandler
	PlayerHandler      *PlayerHandler
	ReserveHandler     *ReserveHandler
	Resources          res.OrderResources
	OrderService       *services.OrderService
	DateHelper         telegram.DateKeyboardHelper
	ListDateHelper     telegram.DateKeyboardHelper
	TimeHelper         telegram.TimeKeyboardHelper
	ActivityHelper     telegram.EnumKeyboardHelper
	MinLevelHelper     telegram.EnumKeyboardHelper
	CourtsHelper       telegram.CountKeyboardHelper
	SetsHelper         telegram.CountKeyboardHelper
	PlayerCountHelper  telegram.CountKeyboardHelper
	JoinCountHelper    telegram.CountKeyboardHelper
	PriceCountHelper   telegram.CountKeyboardHelper
	OrderActionsHelper telegram.ActionsKeyboardHelper
}

func (oh *OrderBotHandler) GetCommands(tuser *telegram.User) (cmds []telegram.BotCommand) {
	cmds = append(cmds, oh.Resources.ListCommand)
	p, err := oh.GetPerson(tuser)
	if err != nil {
		return
	}
	l, err := oh.GetLocation(oh.Resources.Location.Name)

	if err != nil {
		return
	}

	if p.CheckLocationRole(l, "admin") || p.CheckLocationRole(l, "order") {
		cmds = append(cmds, oh.Resources.OrderCommand)
	}
	return
}

func (oh *OrderBotHandler) GetCallbackHandlers() (hlist []telegram.CallbackHandler) {
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderlistdate", Handler: oh.ListDateCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderdate", Handler: oh.StartDateCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordertime", Handler: oh.StartTimeCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderminlevel", Handler: oh.MinLevelCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderactivity", Handler: oh.ActivityCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordercourts", Handler: oh.CourtsCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordersets", Handler: oh.SetsCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderplayers", Handler: oh.MaxPlayersCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderprice", Handler: oh.PriceCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderjoin", Handler: oh.JoinCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderremovepl", Handler: oh.RemovePlayerCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderjoinmult", Handler: oh.JoinMultiCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderarrivetime", Handler: oh.ArriveTimeCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderleave", Handler: oh.LeaveCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordercancel", Handler: oh.CancelCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordercancelcomfirm", Handler: oh.CancelComfirmCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordershow", Handler: oh.ShowCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderlist", Handler: oh.ListOrdersCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderpub", Handler: oh.PublishCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderdesc", Handler: oh.DescriptionCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordercopy", Handler: oh.CopyCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "ordersettings", Handler: oh.ShowCallback})
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderactions", Handler: oh.ShowCallback})
	return
}

func (oh *OrderBotHandler) GetMessageHandler() (hlist []telegram.MessageHandler) {
	hlist = append(hlist, &telegram.CommandHandler{
		Command: oh.Resources.OrderCommand.Command, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return oh.CreateOrder(m, nil)
		}})
	hlist = append(hlist, &telegram.CommandHandler{
		Command: oh.Resources.ListCommand.Command, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return oh.ListOrders(m, nil), nil
		}})
	hlist = append(hlist, &telegram.StateMessageHandler{State: "orderdesc", StateRepository: oh.StateRepository,
		Handler: oh.DescriptionState})
	return
}

func (oh *OrderBotHandler) GetLocation(lname string) (l location.Location, err error) {
	l, err = oh.OrderService.Locations.GetByName(lname)
	if err != nil {
		log.Println(err.Error())
		l, _ = location.NewLocation(lname)
		l, err = oh.OrderService.Locations.Add(l)
	}
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("getting location error: %s", err.Error()),
			AnswerMsg: "Can't get location"}
	}
	return
}

func (oh *OrderBotHandler) CreateOrder(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(msg.From)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	l, err := oh.GetLocation(oh.Resources.Location.Name)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	if !(p.CheckLocationRole(l, "admin") || p.CheckLocationRole(l, "order")) {
		err = telegram.HelperError{
			Msg:       "Command \"*order*\" not permited",
			AnswerMsg: "Command \"order\" not permited"}
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)

	}
	currTime := time.Now()
	stime := time.Date(currTime.Year(), currTime.Month(), currTime.Day(),
		currTime.Hour()+1, 0, 0, 0, currTime.Location())
	etime := stime.Add(time.Duration(time.Hour))

	res, err := oh.OrderService.CreateOrder(reserve.Reserve{
		Person: p, StartTime: stime, EndTime: etime, Location: l}, nil)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("creating order error: %s", err.Error()),
			AnswerMsg: "Can't create order"}
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	var kbd telegram.InlineKeyboardMarkup
	rm := NewReserveMessager(res, nil, oh.Resources)
	rm.SetReserveActions(p, msg.Chat.Id, "ordershow")

	rm.KeyboardHelper.SetData(res.Id.String())
	kbd.InlineKeyboard = rm.KeyboardHelper.GetKeyboard()
	rview := reserve.NewTelegramViewRu(res)
	mr := &telegram.MessageRequest{
		ChatId:      msg.Chat.Id,
		Text:        rview.GetText(),
		ParseMode:   rview.ParseMode,
		ReplyMarkup: kbd}
	result = oh.Bot.SendMessage(mr)
	if chanr != nil {
		chanr <- result
	}
	return result, nil
}

func (oh *OrderBotHandler) ListOrders(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse) {
	var kbd telegram.InlineKeyboardMarkup
	kbd.InlineKeyboard = oh.ListDateHelper.GetKeyboard()
	mr := &telegram.MessageRequest{
		ChatId:      msg.Chat.Id,
		Text:        oh.ListDateHelper.Msg,
		ReplyMarkup: kbd}
	result = oh.Bot.SendMessage(mr)
	if chanr != nil {
		chanr <- result
	}
	return result
}

func (oh *OrderBotHandler) ListOrdersCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	var kbd telegram.InlineKeyboardMarkup
	kbd.InlineKeyboard = oh.ListDateHelper.GetKeyboard()
	mr := &telegram.EditMessageTextRequest{
		ChatId:      cq.Message.Chat.Id,
		MessageId:   cq.Message.MessageId,
		Text:        oh.ListDateHelper.Msg,
		ReplyMarkup: kbd}
	result = oh.Bot.SendMessage(mr)
	return
}

func (oh *OrderBotHandler) ListDateCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	err = oh.ListDateHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rlist, err := oh.OrderService.List(reserve.Reserve{
		StartTime: oh.ListDateHelper.Date,
		EndTime:   oh.ListDateHelper.Date.Add(time.Duration(time.Hour * 24))}, true, nil)

	mr := telegram.EditMessageTextRequest{ChatId: cq.Message.Chat.Id, MessageId: cq.Message.MessageId}
	if len(rlist) == 0 {
		cq.Answer(oh.Bot, oh.Resources.NoReservesAnswer, nil)
		mr.Text = fmt.Sprintf(oh.Resources.NoReservesMessage,
			monday.Format(oh.ListDateHelper.Date, "Monday, 02.01.2006", oh.Resources.Locale))
		mr.ReplyMarkup = telegram.InlineKeyboardMarkup{InlineKeyboard: oh.ListDateHelper.GetKeyboard()}

		result = oh.Bot.SendMessage(&mr)
		return

	}
	ah := telegram.ActionsKeyboardHelper{Columns: 1}
	prefix := "ordershow"
	for _, res := range rlist {
		tgv := reserve.NewTelegramViewRu(res)
		ab := telegram.ActionButton{
			Prefix: prefix, Data: res.Id.String(), Text: tgv.String()}
		ah.Actions = append(ah.Actions, ab)
	}
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "orderlist", Text: oh.Resources.BackBtn})

	mr.Text = oh.Resources.ReservesMessage
	mr.ReplyMarkup = telegram.InlineKeyboardMarkup{InlineKeyboard: ah.GetKeyboard()}
	result = oh.Bot.SendMessage(&mr)

	cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil)
	return
}

func (oh *OrderBotHandler) StartDateCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	dh := oh.DateHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &dh, nil)

	if dh.Action == "set" {
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		dur := res.GetDuration()
		res.StartTime = dh.Date.Add(time.Duration(res.StartTime.Hour()*int(time.Hour) +
			res.StartTime.Minute()*int(time.Minute)))
		res.EndTime = res.StartTime.Add(dur)

		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		rm := NewReserveMessager(res, &dh, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) StartTimeCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	th := oh.TimeHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &th, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if th.Action == "set" {
		dur := res.GetDuration()
		res.StartTime = time.Date(res.StartTime.Year(), res.StartTime.Month(), res.StartTime.Day(),
			th.Time.Hour(), th.Time.Minute(), 0, 0, time.Local)
		res.EndTime = res.StartTime.Add(dur)
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		th.Step = 30
		rm := NewReserveMessager(res, &th, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) ArriveTimeCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	th := oh.TimeHelper
	th.Prefix = "orderarrivetime"
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &th, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if th.Action == "set" {
		pl := res.GetPlayerByTelegramId(cq.From.Id)
		pl.ArriveTime = th.Time
		res.JoinPlayer(pl)
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		th.StartHour = res.StartTime.Hour()
		th.EndHour = res.EndTime.Hour() - 1
		th.Step = 15
		rm := NewReserveMessager(res, &th, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) SetsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.SetsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.EndTime = res.StartTime.Add(time.Duration(time.Hour * time.Duration(ch.Count)))
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) CopyCallback(cq *telegram.CallbackQuery) (resp telegram.MessageResponse, err error) {
	p, resp, err := oh.GetPersonCq(cq)
	if err != nil {
		return
	}

	ch := oh.OrderActionsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err = oh.OrderService.Reserves.Add(res.Copy())
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("copping order error: %s", err.Error()),
			AnswerMsg: "Can't copy order"}
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rm := NewReserveMessager(res, nil, oh.Resources)
	rm.SetReserveActions(p, cq.Message.Chat.Id, "ordercopy")
	mr := rm.GetEditMR(cq.Message.Chat.Id)
	cq.Message.EditText(oh.Bot, oh.Resources.CopyMessage, &mr)
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) ShowCallback(cq *telegram.CallbackQuery) (resp telegram.MessageResponse, err error) {
	ch := oh.OrderActionsHelper
	p, resp, err := oh.GetPersonCq(cq)
	if err != nil {
		return
	}
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rm := NewReserveMessager(res, nil, oh.Resources)
	rm.SetReserveActions(p, cq.Message.Chat.Id, ch.Action)
	mr := rm.GetEditMR(cq.Message.Chat.Id)
	cq.Message.EditText(oh.Bot, "", &mr)
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) PublishCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.OrderActionsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	st := telegram.State{
		ChatId: res.Location.ChatId,
		State:  "ordershow",
		Data:   res.Id.String(),
	}
	rm := NewReserveMessager(res, nil, oh.Resources)
	rm.SetReserveActions(person.Person{}, st.ChatId, st.State)
	mr := rm.GetMR(st.ChatId)
	resp := oh.Bot.SendMessage(&mr)
	st.MessageId = resp.Result.MessageId
	oh.StateRepository.Set(st)

	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) MinLevelCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.MinLevelHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		lvl, err := strconv.Atoi(ch.Choice)
		res.MinLevel = lvl
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) ActivityCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.ActivityHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		act, err := strconv.Atoi(ch.Choice)
		res.Activity = reserve.Activity(act)
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) CourtsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.CourtsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.CourtCount = ch.Count
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		ch.Max = res.Location.CourtCount
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) DescriptionState(msg *telegram.Message, state telegram.State) (resp telegram.MessageResponse, err error) {
	res, err := oh.ReserveHandler.GetDataReserve(state.Data, nil, nil)
	if err != nil {
		state.State = "ordershow"
		oh.StateRepository.Set(state)
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	res.Description = msg.Text
	state.State = "ordershow"
	oh.StateRepository.Set(state)
	upmsg := *msg
	upmsg.MessageId = state.MessageId
	if resp, err = oh.ReserveHandler.UpdateReserveMsg(res, &upmsg, state.MessageId); err != nil {
		return
	}
	return oh.Bot.SendMessage(msg.CreateMessageRequest(oh.Resources.Description.DoneMessage, nil)), nil
}

func (oh *OrderBotHandler) DescriptionCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	ah.Msg = oh.Resources.Description.Message
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ah, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	oh.StateRepository.Set(telegram.State{
		State:     "orderdesc",
		ChatId:    cq.Message.Chat.Id,
		Data:      res.Id.String(),
		MessageId: cq.Message.MessageId,
	})
	oh.Bot.SendMessage(cq.Message.CreateMessageRequest(oh.Resources.Description.Message, nil))
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) MaxPlayersCallback(cq *telegram.CallbackQuery) (resp telegram.MessageResponse, err error) {
	ch := oh.PlayerCountHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		if ch.Count < res.PlayerCount(uuid.Nil) {
			err = telegram.HelperError{Msg: "Max player count error.", AnswerMsg: oh.Resources.MaxPlayer.CountError}
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		res.MaxPlayers = ch.Count
		oh.StateRepository.Clear(telegram.State{ChatId: cq.Message.Chat.Id, MessageId: cq.Message.MessageId})
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		ch.Max = res.CourtCount * 12
		oh.StateRepository.Set(telegram.State{
			State:     "orderplayers",
			ChatId:    cq.Message.Chat.Id,
			Data:      res.Id.String(),
			MessageId: cq.Message.MessageId,
		})
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) PriceCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.PriceCountHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.Price = ch.Count
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		rm := NewReserveMessager(res, &ch, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) JoinPlayer(cq *telegram.CallbackQuery, Data string, count int) (result telegram.MessageResponse, err error) {
	res, err := oh.ReserveHandler.GetDataReserve(Data, nil, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	oh.PlayerHandler.JoinPlayer(cq, &res, count)
	return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", true)
}

func (oh *OrderBotHandler) JoinCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	return oh.JoinPlayer(cq, ah.Data, 1)
}

func (oh *OrderBotHandler) RemovePlayerCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ph := telegram.NewEnumKeyboardHelper(oh.Resources.Activity.Message, "orderremovepl", nil)
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ph, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ph.Action == "set" {
		tid, err := strconv.Atoi(ph.Choice)
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		pl := res.GetPlayerByTelegramId(tid)
		pl.Count = 0
		res.JoinPlayer(pl)
		return oh.ReserveHandler.UpdateReserveCQ(res, cq, "orderactions", false)
	} else {
		pllist := []telegram.EnumItem{}
		for _, pl := range res.Players {
			pllist = append(pllist, telegram.EnumItem{Id: strconv.Itoa(pl.TelegramId), Item: pl.String()})
		}
		ph := telegram.NewEnumKeyboardHelper(oh.Resources.Activity.Message, "orderremovepl", pllist)

		ab := telegram.ActionButton{
			Prefix: "orderactions", Data: res.Id.String(), Text: ""}
		ph.BackData = telegram.ActionsKeyboardHelper{}.GetBtnData(ab)
		rm := NewReserveMessager(res, &ph, oh.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) JoinMultiCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	ch := oh.JoinCountHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ch, nil)

	if ch.Action == "set" {
		return oh.JoinPlayer(cq, ch.Data, ch.Count)
	} else {
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}

		ch.Min = 1
		if cq.Message.Chat.Id > 0 {
			ch.Max = res.MaxPlayers - res.PlayerCount(p.Id)
			if !res.HasPlayerByTelegramId(p.TelegramId) || res.PlayerInReserve(p.Id) {
				ch.Max = res.MaxPlayers
			} else if ch.Max <= res.GetPlayer(p.Id).Count {
				ch.Max = res.GetPlayer(p.Id).Count
			}
		} else {
			ch.Max = res.MaxPlayers - res.PlayerCount(uuid.Nil)
		}
		if ch.Max > 1 {
			rm := NewReserveMessager(res, &ch, oh.Resources)
			mr := rm.GetEditMR(cq.Message.Chat.Id)
			cq.Message.EditText(oh.Bot, "", &mr)
		}
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) LeaveCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	return oh.JoinPlayer(cq, ah.Data, 0)
}

func (oh *OrderBotHandler) CancelCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ah, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	ah.Columns = 2
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "orderleave", Text: oh.Resources.JoinPlayer.LeaveButton})
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "ordercancelcomfirm", Text: oh.Resources.Cancel.Confirm})
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "ordershow", Text: oh.Resources.Cancel.Abort})
	rm := NewReserveMessager(res, &ah, oh.Resources)
	mr := rm.GetEditMR(cq.Message.Chat.Id)
	mr.Text += oh.Resources.Cancel.Message
	cq.Message.EditText(oh.Bot, "", &mr)

	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) CancelComfirmCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	res, err := oh.ReserveHandler.GetDataReserve(cq.Data, &ah, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.Canceled = true
	rm := NewReserveMessager(res, nil, oh.Resources)
	oh.PlayerHandler.NotifyPlayers(res, &rm, cq.From.Id, "notify_cancel")
	return oh.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
}
