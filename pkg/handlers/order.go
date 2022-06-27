package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/goodsign/monday"
	"github.com/google/uuid"
)

func NewOrderHandler(tb *telegram.Bot, os *services.OrderService, rl OrderResourceLoader) (oh OrderBotHandler) {
	oh = OrderBotHandler{OrderService: os}
	oh.Bot = tb
	oh.PersonService = os.PersonService
	oh.Resources = rl.GetResource()

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
	Resources          OrderResources
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
	hlist = append(hlist, &telegram.PrefixCallbackHandler{Prefix: "orderjoinmult", Handler: oh.JoinMultiCallback})
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
	kh := oh.GetReserveActions(res, p, msg.Chat.Id, "ordershow")
	kh.SetData(res.Id.String())
	kbd.InlineKeyboard = kh.GetKeyboard()
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

func (oh *OrderBotHandler) GetDataReserve(data string,
	rchan chan services.ReserveResult) (r reserve.Reserve, err error) {
	var id uuid.UUID
	id, err = uuid.Parse(data)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order getting reserve error: %s", err.Error()),
			AnswerMsg: "Parse reserve id error"}

	} else {
		r, err = oh.OrderService.Reserves.Get(id)
		if err != nil {
			err = telegram.HelperError{
				Msg:       fmt.Sprintf("Getting reserve error: %s", err.Error()),
				AnswerMsg: "Getting reserve error"}
		}
	}

	if rchan != nil {
		rchan <- services.ReserveResult{Reserve: r, Err: err}
	}
	return
}

func (oh *OrderBotHandler) GetReserveActions(res reserve.Reserve, p person.Person, chid int, state string) (h telegram.KeyboardHelper) {
	ah := telegram.ActionsKeyboardHelper{Data: res.Id.String()}
	if res.Canceled {
		return &ah
	}
	ah.Columns = 2
	if chid == p.TelegramId {
		if res.Person.TelegramId == p.TelegramId || p.CheckLocationRole(res.Location, "admin") {
			if state == "ordershow" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderdate", Text: oh.Resources.DateTime.DateButton})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordertime", Text: oh.Resources.DateTime.TimeButton})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordersets", Text: oh.Resources.Set.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderdesc", Text: oh.Resources.Description.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordersettings", Text: oh.Resources.SettingsBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderactions", Text: oh.Resources.ActionsBtn})
			} else if state == "ordersettings" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderactivity", Text: oh.Resources.Activity.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercourts", Text: oh.Resources.Court.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderminlevel", Text: oh.Resources.Level.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderplayers", Text: oh.Resources.MaxPlayer.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderprice", Text: oh.Resources.Price.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordershow", Text: oh.Resources.BackBtn})
			} else if state == "orderactions" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercancel", Text: oh.Resources.Cancel.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercopy", Text: oh.Resources.CopyBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderpub", Text: oh.Resources.PublishBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordershow", Text: oh.Resources.BackBtn})
			}
		}
	}
	if res.Ordered() && state == "ordershow" {
		if chid <= 0 || !res.HasPlayerByTelegramId(p.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderjoin", Text: oh.Resources.JoinPlayer.Button})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderjoinmult", Text: oh.Resources.JoinPlayer.MultiButton})
		if chid <= 0 || res.HasPlayerByTelegramId(p.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderleave", Text: oh.Resources.JoinPlayer.LeaveButton})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordershow", Text: oh.Resources.RefreshBtn})
	}

	return &ah
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
	err = dh.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(dh.GetData(), nil)

	if dh.Action == "set" {
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		dur := res.GetDuration()
		res.StartTime = dh.Date.Add(time.Duration(res.StartTime.Hour()*int(time.Hour) +
			res.StartTime.Minute()*int(time.Minute)))
		res.EndTime = res.StartTime.Add(dur)

		return oh.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		mr := oh.GetReserveEditMR(res, &dh)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) StartTimeCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	th := oh.TimeHelper
	err = th.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(th.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if th.Action == "set" {
		dur := res.GetDuration()
		res.StartTime = time.Date(res.StartTime.Year(), res.StartTime.Month(), res.StartTime.Day(),
			th.Time.Hour(), 0, 0, 0, time.Local)
		res.EndTime = res.StartTime.Add(dur)
		return oh.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		mr := oh.GetReserveEditMR(res, &th)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) SetsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.SetsHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.EndTime = res.StartTime.Add(time.Duration(time.Hour * time.Duration(ch.Count)))
		return oh.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
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
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(ch.Data, nil)
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

	kbd := oh.GetReserveActions(res, p, cq.Message.Chat.Id, "ordercopy")
	mr := oh.GetReserveEditMR(res, kbd)
	mr.ChatId = cq.Message.Chat.Id
	cq.Message.EditText(oh.Bot, oh.Resources.CopyMessage, &mr)
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) ShowCallback(cq *telegram.CallbackQuery) (resp telegram.MessageResponse, err error) {
	ch := oh.OrderActionsHelper
	p, resp, err := oh.GetPersonCq(cq)
	if err != nil {
		return
	}

	if err = ch.Parse(cq.Data); err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	kbd := oh.GetReserveActions(res, p, cq.Message.Chat.Id, ch.Action)
	mr := oh.GetReserveEditMR(res, kbd)
	mr.ChatId = cq.Message.Chat.Id
	cq.Message.EditText(oh.Bot, "", &mr)
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) PublishCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.OrderActionsHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	st := telegram.State{
		ChatId: res.Location.ChatId,
		State:  "ordershow",
		Data:   res.Id.String(),
	}
	kbd := oh.GetReserveActions(res, person.Person{}, st.ChatId, st.State)
	mr := oh.GetReserveMR(res, kbd)
	mr.ChatId = st.ChatId
	resp := oh.Bot.SendMessage(&mr)
	st.MessageId = resp.Result.MessageId
	oh.StateRepository.Set(st)

	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) MinLevelCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.MinLevelHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		lvl, err := strconv.Atoi(ch.Choice)
		res.MinLevel = lvl
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		return oh.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) ActivityCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.ActivityHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		act, err := strconv.Atoi(ch.Choice)
		res.Activity = act
		if err != nil {
			return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		return oh.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) CourtsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.CourtsHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.CourtCount = ch.Count
		return oh.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		ch.Max = res.Location.CourtCount
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) DescriptionState(msg *telegram.Message, state telegram.State) (resp telegram.MessageResponse, err error) {
	res, err := oh.GetDataReserve(state.Data, nil)
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
	if resp, err = oh.UpdateReserveMsg(res, &upmsg, state.MessageId); err != nil {
		return
	}
	return oh.Bot.SendMessage(msg.CreateMessageRequest(oh.Resources.Description.DoneMessage, nil)), nil
}

func (oh *OrderBotHandler) DescriptionCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	ah.Msg = oh.Resources.Description.Message
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ah.Data, nil)
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
	if err = ch.Parse(cq.Data); err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
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
		return oh.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		ch.Max = res.CourtCount * 12
		oh.StateRepository.Set(telegram.State{
			State:     "orderplayers",
			ChatId:    cq.Message.Chat.Id,
			Data:      res.Id.String(),
			MessageId: cq.Message.MessageId,
		})
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) PriceCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.PriceCountHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.Price = ch.Count
		return oh.UpdateReserveCQ(res, cq, "ordersettings", false)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) JoinCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	return oh.JoinPlayer(cq, ah.Data, 1)
}

func (oh *OrderBotHandler) JoinMultiCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	ch := oh.JoinCountHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if ch.Action == "set" {
		return oh.JoinPlayer(cq, ch.Data, ch.Count)
	} else {
		ch.Min = 1
		ch.Max = res.MaxPlayers - res.PlayerCount(p.Id)
		if ch.Max <= res.GetPlayer(p.Id).Count {
			ch.Max = res.GetPlayer(p.Id).Count
		} else if res.GetPlayer(p.Id).Count == 0 {
			ch.Max = res.MaxPlayers
		}
		var mr telegram.EditMessageTextRequest
		if ch.Max > 1 {
			mr = oh.GetReserveEditMR(res, &ch)
		} else {
			oh.GetReserveActions(res, p, cq.Message.Chat.Id, "orderjoin")
		}
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) JoinPlayer(cq *telegram.CallbackQuery, data string, count int) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.JoinPlayer(person.Player{Person: p, Count: count})
	return oh.UpdateReserveCQ(res, cq, "ordershow", true)
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
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(ah.Data, nil)
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
	mr := oh.GetReserveEditMR(res, &ah)
	mr.ChatId = cq.Message.Chat.Id
	mr.Text += oh.Resources.Cancel.Message
	cq.Message.EditText(oh.Bot, "", &mr)

	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) CancelComfirmCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ah := oh.OrderActionsHelper
	err = ah.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(ah.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.Canceled = true
	for _, pl := range res.Players {
		if pl.Person.TelegramId != cq.From.Id {
			p, _ := oh.OrderService.PersonService.GetByTelegramId(pl.Person.TelegramId)
			if p.Settings["notify_cancel"] != "off" {
				mr := oh.GetReserveMR(res, nil)
				mr.ChatId = pl.Person.TelegramId
				oh.Bot.SendMessage(&mr)
			}
		}
	}
	return oh.UpdateReserveCQ(res, cq, "ordershow", false)
}

func (oh *OrderBotHandler) NotifyPlayers(res reserve.Reserve, id int) {
	for _, pl := range res.Players {
		if pl.Person.TelegramId != id {
			p, _ := oh.OrderService.PersonService.GetByTelegramId(pl.Person.TelegramId)
			if param, ok := p.Settings["notify"]; ok && param == "on" {
				mr := oh.GetReserveMR(res, nil)
				mr.ChatId = pl.Person.TelegramId
				oh.Bot.SendMessage(&mr)
				return
			}
		}
	}
}

func (oh *OrderBotHandler) UpdateReserveCQ(res reserve.Reserve, cq *telegram.CallbackQuery, state string, renew bool) (resp telegram.MessageResponse, err error) {
	st := telegram.State{
		ChatId:    cq.Message.Chat.Id,
		Data:      res.Id.String(),
		State:     state,
		MessageId: cq.Message.MessageId,
	}
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err = oh.UpdateReserve(res)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	msg := cq.Message
	if renew && st.ChatId < 0 {
		mr := oh.GetReserveMR(res, oh.GetReserveActions(res, p, st.ChatId, st.State))
		mr.DisableNotification = true
		resp = msg.SendMessage(oh.Bot, "", &mr)
		msg.DeleteMessage(oh.Bot)
		oh.StateRepository.Clear(telegram.State{ChatId: st.ChatId, MessageId: msg.MessageId})
		msg = &resp.Result
		st.MessageId = resp.Result.MessageId
	} else {
		mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, p, st.ChatId, st.State))
		cq.Message.EditText(oh.Bot, "", &mr)
	}
	oh.StateRepository.Set(st)

	oh.UpdateReserveMessages(res, msg, true)
	oh.NotifyPlayers(res, cq.From.Id)
	resp = cq.Answer(oh.Bot, "Ok", nil)

	return resp, nil
}

func (oh *OrderBotHandler) UpdateReserveMsg(res reserve.Reserve, msg *telegram.Message, mid int) (resp telegram.MessageResponse, err error) {
	st := telegram.State{Data: res.Id.String(), State: "ordershow"}
	p, err := oh.GetPerson(msg.From)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	res, err = oh.UpdateReserve(res)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, p, msg.Chat.Id, st.State))
	if mid > 0 {
		mr.ChatId = msg.Chat.Id
		mr.MessageId = mid
	}
	resp = oh.Bot.SendMessage(&mr)
	oh.UpdateReserveMessages(res, msg, true)
	st.ChatId = mr.ChatId.(int)
	st.MessageId = resp.Result.MessageId
	oh.StateRepository.Set(st)
	oh.NotifyPlayers(res, msg.From.Id)

	return
}

func (oh *OrderBotHandler) UpdateReserveMessages(res reserve.Reserve, msg *telegram.Message, renew bool) {
	slist, _ := oh.StateRepository.GetByData(res.Id.String())
	for _, st := range slist {
		if msg.MessageId == st.MessageId && msg.Chat.Id == st.ChatId {
			continue
		}
		p, _ := oh.GetPerson(&telegram.User{Id: st.ChatId})
		if renew && st.ChatId < 0 {
			mr := oh.GetReserveMR(res, oh.GetReserveActions(res, p, st.ChatId, st.State))
			mr.ChatId = st.ChatId
			mr.DisableNotification = true
			resp := oh.Bot.SendMessage(&mr)
			oh.Bot.SendMessage(&telegram.DeleteMessageRequest{ChatId: st.ChatId, MessageId: st.MessageId})
			oh.StateRepository.Clear(st)
			st.MessageId = resp.Result.MessageId
			oh.StateRepository.Set(st)
		} else {
			mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, p, st.ChatId, st.State))
			mr.ChatId = st.ChatId
			mr.MessageId = st.MessageId
			oh.Bot.SendMessage(&mr)
		}
	}
}

func (oh *OrderBotHandler) GetReserveEditMR(res reserve.Reserve, kh telegram.KeyboardHelper) (mer telegram.EditMessageTextRequest) {
	mr := oh.GetReserveMR(res, kh)
	return telegram.EditMessageTextRequest{Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (oh *OrderBotHandler) GetReserveMR(res reserve.Reserve, kh telegram.KeyboardHelper) (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	if kh != nil {
		kh.SetData(res.Id.String())
		kbd.InlineKeyboard = append(kbd.InlineKeyboard, kh.GetKeyboard()...)
		kbdText = "\n*" + kh.GetText() + "* "
	}

	rview := reserve.NewTelegramViewRu(res)
	mtxt := fmt.Sprintf("%s\n%s", rview.GetText(), kbdText)
	if len(kbd.InlineKeyboard) > 0 {
		return telegram.MessageRequest{Text: mtxt, ParseMode: rview.ParseMode, ReplyMarkup: kbd}
	}
	return telegram.MessageRequest{Text: mtxt, ParseMode: rview.ParseMode}
}

func (oh *OrderBotHandler) GetReserve(id uuid.UUID) (result reserve.Reserve, err error) {
	result, err = oh.OrderService.Reserves.Get(id)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order update can't get reserve %s error: %s", id, err.Error()),
			AnswerMsg: "Can't get reserve"}
	}
	return
}

func (oh *OrderBotHandler) UpdateReserve(res reserve.Reserve) (result reserve.Reserve, err error) {
	err = oh.OrderService.Reserves.Update(res)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order update can't update reserve %s error: %s", res.Id, err.Error()),
			AnswerMsg: "Can't update reserve"}
		return
	}
	result, err = oh.GetReserve(res.Id)
	return
}
