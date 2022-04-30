package handlers

import (
	"fmt"
	"log"
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

func NewOrderHandler(tb *telegram.Bot, os *services.OrderService) (oh OrderBotHandler) {
	oh = OrderBotHandler{Bot: tb, OrderService: os}
	oh.OrderCommand = "order"
	oh.ListCommand = "list"
	oh.DateHelper = telegram.NewDateKeyboardHelper("Выбери дату:", "orderdate")
	oh.ListDateHelper = telegram.NewDateKeyboardHelper("Выбери дату:", "orderlistdate")
	oh.TimeHelper = telegram.NewTimeKeyboardHelper("Выбери время:", "ordertime")
	oh.CourtsHelper = telegram.NewCountKeyboardHelper("❓Сколько нужно кортов❓", "ordercourts", 1, 4)
	oh.SetsHelper = telegram.NewCountKeyboardHelper("❓Количество часов❓", "ordersets", 1, 4)
	oh.PlayerCountHelper = telegram.NewCountKeyboardHelper("❓Максимальное количество игроков❓", "orderplayers", 4, 32)
	oh.JoinCountHelper = telegram.NewCountKeyboardHelper("❓Сколько игроков записать❓", "orderjoincount", 1, 32)
	oh.OrderActionsHelper.Columns = 2
	oh.OrderActionsHelper.Actions = append(oh.OrderActionsHelper.Actions, telegram.ActionButton{Prefix: "orderjoin", Text: "Присоедениться"})
	oh.OrderActionsHelper.Actions = append(oh.OrderActionsHelper.Actions, telegram.ActionButton{Prefix: "orderleave", Text: "Не пойду"})

	return
}

type OrderBotHandler struct {
	OrderCommand       string
	ListCommand        string
	Bot                *telegram.Bot
	OrderService       *services.OrderService
	DateHelper         telegram.DateKeyboardHelper
	ListDateHelper     telegram.DateKeyboardHelper
	TimeHelper         telegram.TimeKeyboardHelper
	CourtsHelper       telegram.CountKeyboardHelper
	SetsHelper         telegram.CountKeyboardHelper
	PlayerCountHelper  telegram.CountKeyboardHelper
	JoinCountHelper    telegram.CountKeyboardHelper
	OrderActionsHelper telegram.ActionsKeyboardHelper
	MessageHandlers    []telegram.MessageHandler
	CallbackHandlers   []telegram.CallbackHandler
}

func (oh *OrderBotHandler) ProceedCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	if len(oh.CallbackHandlers) == 0 {
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderlistdate", Handler: oh.ListDateCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderdate", Handler: oh.StartDateCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordertime", Handler: oh.StartTimeCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercourts", Handler: oh.CourtsCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordersets", Handler: oh.SetsCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderplayers", Handler: oh.MaxPlayersCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderjoin", Handler: oh.JoinCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderjoincount", Handler: oh.JoinCountCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderleave", Handler: oh.LeaveCallback})
	}
	for _, handler := range oh.CallbackHandlers {
		result, err = handler.ProceedCallback(cq)
	}
	return
}

func (oh *OrderBotHandler) ProceedMessage(msg *telegram.Message) (result telegram.MessageResponse, err error) {
	if len(oh.MessageHandlers) == 0 {
		order_cmd := telegram.CommandHandler{Command: oh.OrderCommand, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return oh.CreateOrder(m, nil)
		}}
		oh.MessageHandlers = append(oh.MessageHandlers, &order_cmd)
		list_cmd := telegram.CommandHandler{Command: oh.ListCommand, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
			return oh.ListOrders(m, nil), nil
		}}
		oh.MessageHandlers = append(oh.MessageHandlers, &list_cmd)
	}
	for _, handler := range oh.MessageHandlers {
		result, err = handler.ProceedMessage(msg)
	}
	return
}

func (oh *OrderBotHandler) SendCallbackError(cq *telegram.CallbackQuery, cq_err telegram.HelperError, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	log.Println(cq_err.Error())
	result = cq.Answer(oh.Bot, cq_err.AnswerMsg, nil)
	if chanr != nil {
		chanr <- result
	}
	return result, cq_err
}

func (oh *OrderBotHandler) SendNessageError(msg *telegram.Message, m_err telegram.HelperError, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	log.Println(m_err.Error())
	result = msg.Reply(oh.Bot, m_err.AnswerMsg, nil)
	if chanr != nil {
		chanr <- result
	}
	return result, m_err
}

func (oh *OrderBotHandler) GetPerson(tuser *telegram.User) (p person.Person, err error) {
	p, err = oh.OrderService.Persons.GetByTelegramId(tuser.Id)
	if err != nil {
		log.Println(err.Error())
		p, _ = person.NewPerson(tuser.FirstName)
		p.TelegramId = tuser.Id
		p.Lastname = tuser.LastName
		p, err = oh.OrderService.Persons.Add(p)
	}
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("getting person error: %s", err.Error()),
			AnswerMsg: "Can't get person"}
	}
	return
}

func (oh *OrderBotHandler) CreateOrder(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(msg.From)
	if err != nil {
		return oh.SendNessageError(msg, err.(telegram.HelperError), nil)
	}

	res := oh.OrderService.CreateOrder(reserve.Reserve{Person: p}, nil)
	if res.Err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("creating order error: %s", err.Error()),
			AnswerMsg: "Can't create order"}
		return oh.SendNessageError(msg, err.(telegram.HelperError), nil)
	}

	var kbd telegram.InlineKeyboardMarkup
	oh.DateHelper.Data = res.Reserve.Id.String()
	kbd.InlineKeyboard = oh.DateHelper.GetKeyboard()
	mr := &telegram.MessageRequest{
		ChatId:      msg.Chat.Id,
		Text:        "Давай для начала выберем дату.",
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

func (oh *OrderBotHandler) GetDataReserve(data string,
	rchan chan services.ReserveResult) (r reserve.Reserve, err error) {
	var id uuid.UUID
	id, err = uuid.Parse(data)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order CourtsCallback getting reserve error: %s", err.Error()),
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

func (oh *OrderBotHandler) ListDateCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	sdate, _, err := oh.ListDateHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rlist := oh.OrderService.List(reserve.Reserve{
		StartTime: sdate,
		EndTime:   sdate.Add(time.Duration(time.Hour * 24)),
		Ordered:   true}, nil)

	for _, res := range rlist.Reserves {
		kh := oh.OrderActionsHelper.SetData(res.Id.String())
		mr := oh.GetReserveMR(res, kh)
		mr.ChatId = cq.Message.Chat.Id
		result = oh.Bot.SendMessage(&mr)
	}
	cq.Answer(oh.Bot, "Ok", nil)
	return
}

func (oh *OrderBotHandler) StartDateCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	sdate, data, err := oh.DateHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.StartTime = sdate.Add(time.Duration(res.StartTime.Hour()*int(time.Hour) +
		res.StartTime.Minute()*int(time.Minute)))

	return oh.UpdateReserveCQ(res, cq, oh.TimeHelper)
}

func (oh *OrderBotHandler) StartTimeCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	stime, data, err := oh.TimeHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.StartTime = time.Date(res.StartTime.Year(), res.StartTime.Month(), res.StartTime.Day(), stime.Hour(), 0, 0, 0, time.Local)
	return oh.UpdateReserveCQ(res, cq, oh.SetsHelper)
}

func (oh *OrderBotHandler) SetsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	count, data, err := oh.SetsHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.EndTime = res.StartTime.Add(time.Duration(time.Hour * time.Duration(count)))
	return oh.UpdateReserveCQ(res, cq, oh.CourtsHelper)
}

func (oh *OrderBotHandler) CourtsCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	count, data, err := oh.CourtsHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.CourtCount = count
	return oh.UpdateReserveCQ(res, cq, oh.PlayerCountHelper)
}

func (oh *OrderBotHandler) MaxPlayersCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	count, data, err := oh.PlayerCountHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res.MaxPlayers = count
	res.Ordered = true
	return oh.UpdateReserveCQ(res, cq, oh.OrderActionsHelper)
}

func (oh *OrderBotHandler) JoinCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	data, err := oh.OrderActionsHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	pl_count := 0
	for pl_id, pl := range res.Players {
		if pl_id != p.Id {
			pl_count += pl.Count
		}
	}
	oh.JoinCountHelper.Max = res.MaxPlayers - pl_count
	return oh.UpdateReserveCQ(res, cq, oh.JoinCountHelper)
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
	res.Players[p.Id] = reserve.Player{Person: p, Count: count}
	return oh.UpdateReserveCQ(res, cq, oh.OrderActionsHelper)
}

func (oh *OrderBotHandler) JoinCountCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	count, data, err := oh.JoinCountHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	return oh.JoinPlayer(cq, data, count)
}

func (oh *OrderBotHandler) LeaveCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	data, err := oh.OrderActionsHelper.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	return oh.JoinPlayer(cq, data, 0)
}

func (oh *OrderBotHandler) UpdateReserveCQ(res reserve.Reserve, cq *telegram.CallbackQuery,
	kh telegram.KeyboardHelper) (telegram.MessageResponse, error) {

	res, err := oh.UpdateReserve(res)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	mr := oh.GetReserveEditMR(res, kh)
	cq.Message.EditText(oh.Bot, "", &mr)

	return cq.Answer(oh.Bot, "Ok", nil), nil
}

func (oh *OrderBotHandler) GetReserveEditMR(res reserve.Reserve, kh telegram.KeyboardHelper) (mer telegram.EditMessageTextRequest) {
	mr := oh.GetReserveMR(res, kh)
	return telegram.EditMessageTextRequest{Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (oh *OrderBotHandler) GetReserveMR(res reserve.Reserve, kh telegram.KeyboardHelper) (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	if kh != nil {
		kh = kh.SetData(res.Id.String())
		kbd.InlineKeyboard = append(kbd.InlineKeyboard, kh.GetKeyboard()...)
		kbdText = "\n*" + kh.GetText() + "* "
	}

	rview := reserve.NewTelegramViewRu(res)
	return telegram.MessageRequest{
		Text: fmt.Sprintf("%s\n%s", rview.GetText(), kbdText), ParseMode: rview.ParseMode, ReplyMarkup: kbd}
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
