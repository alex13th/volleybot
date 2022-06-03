package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type DateTimeResources struct {
	DateMessage string
	DateButton  string
	TimeMessage string
	TimeButton  string
}

type CourtResources struct {
	Message    string
	Button     string
	Min        int
	Max        int
	MaxPlayers int
}

type PlayerLevelResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type SetResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type MaxPlayerResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type JoinPlayerResources struct {
	Message     string
	Button      string
	LeaveButton string
}

type PriceResources struct {
	Message string
	Button  string
	Min     int
	Max     int
	Step    int
}

type CancelResources struct {
	Message string
	Button  string
	Confirm string
	Abort   string
}

type OrderResources struct {
	DateTime          DateTimeResources
	Court             CourtResources
	Level             PlayerLevelResources
	Set               SetResources
	MaxPlayer         MaxPlayerResources
	JoinPlayer        JoinPlayerResources
	Price             PriceResources
	Cancel            CancelResources
	NoReservesMessage string
	NoReservesAnswer  string
	OkAnswer          string
}

type OrderResourceLoader interface {
	GetResource() OrderResources
}

type DefaultResourceLoader struct{}

func (rl DefaultResourceLoader) GetResource() (or OrderResources) {
	or.DateTime.DateMessage = "❓Какая дата❓"
	or.DateTime.DateButton = "📆 Дата"
	or.DateTime.TimeMessage = "❓В какое время❓"
	or.DateTime.TimeButton = "⏰ Время"
	or.Level.Message = "❓Какой минимальный уровень игроков❓"
	or.Level.Button = "💪 Уровень"
	or.Set.Message = "❓Количество часов❓"
	or.Set.Button = "⏱ Кол-во часов"
	or.Set.Max = 12
	or.Court.Message = "❓Сколько нужно кортов❓"
	or.Court.Button = "🏐 Площадки"
	or.Court.Max = 12
	or.Court.MaxPlayers = 6
	or.MaxPlayer.Message = "❓Максимальное количество игроков❓"
	or.MaxPlayer.Button = "👫 Мест"
	or.MaxPlayer.Min = 4
	or.MaxPlayer.Max = or.Court.Max * or.Court.MaxPlayers
	or.JoinPlayer.Message = "❓Сколько игроков записать❓"
	or.JoinPlayer.Button = "😀 Хочу"
	or.JoinPlayer.LeaveButton = "😞 Не хочу"
	or.Price.Message = "❓Почем будет поигать❓"
	or.Price.Button = "💳 Стоимость"
	or.Price.Min = 0
	or.Price.Max = 1200
	or.Price.Step = 200
	or.Cancel.Button = "💥Отменить"
	or.Cancel.Message = "\n🧨*ВНИМАНИЕ!!!*🧨\nИгра будет отменена для всех участников. Если есть желание только выписаться, лучше воспользоваться опцией кнопкой \"😞 Не хочу\""
	or.Cancel.Confirm = "🧨 Уверен"
	or.Cancel.Abort = "👌 Передумал"
	or.NoReservesMessage = "На эту дату нет доступных записей."
	or.NoReservesAnswer = "Резервы отсутствуют"
	or.OkAnswer = "Ok"

	return
}

func NewOrderHandler(tb *telegram.Bot, os *services.OrderService, rl OrderResourceLoader) (oh OrderBotHandler) {
	oh = OrderBotHandler{Bot: tb, OrderService: os}
	oh.OrderCommand = "order"
	oh.ListCommand = "list"
	oh.Resources = rl.GetResource()

	oh.DateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderdate")
	oh.ListDateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderlistdate")
	oh.TimeHelper = telegram.NewTimeKeyboardHelper(oh.Resources.DateTime.TimeMessage, "ordertime")

	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		oh.Resources.Price.Step = 200

		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i), Item: reserve.PlayerLevel(i)})
	}
	oh.MinLevelHelper = telegram.NewEnumKeyboardHelper(oh.Resources.Level.Message, "orderminlevel", levels)

	oh.CourtsHelper = telegram.NewCountKeyboardHelper(oh.Resources.Court.Message, "ordercourts", 1, oh.Resources.Court.Max)
	oh.SetsHelper = telegram.NewCountKeyboardHelper(oh.Resources.Set.Message, "ordersets", 1, oh.Resources.Court.Max)
	oh.PlayerCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.MaxPlayer.Message, "orderplayers", oh.Resources.MaxPlayer.Min, oh.Resources.MaxPlayer.Max)
	oh.JoinCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.JoinPlayer.Message, "orderjoin", oh.Resources.MaxPlayer.Min, oh.Resources.MaxPlayer.Max)
	oh.PriceCountHelper = telegram.NewCountKeyboardHelper(
		oh.Resources.Price.Message, "orderprice", oh.Resources.Price.Min, oh.Resources.Price.Max)
	oh.PriceCountHelper.Step = oh.Resources.Price.Step

	return
}

type OrderBotHandler struct {
	Resources          OrderResources
	OrderCommand       string
	ListCommand        string
	Bot                *telegram.Bot
	OrderService       *services.OrderService
	DateHelper         telegram.DateKeyboardHelper
	ListDateHelper     telegram.DateKeyboardHelper
	TimeHelper         telegram.TimeKeyboardHelper
	MinLevelHelper     telegram.EnumKeyboardHelper
	CourtsHelper       telegram.CountKeyboardHelper
	SetsHelper         telegram.CountKeyboardHelper
	PlayerCountHelper  telegram.CountKeyboardHelper
	JoinCountHelper    telegram.CountKeyboardHelper
	PriceCountHelper   telegram.CountKeyboardHelper
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
			&telegram.PrefixCallbackHandler{Prefix: "orderminlevel", Handler: oh.MinLevelCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercourts", Handler: oh.CourtsCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordersets", Handler: oh.SetsCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderplayers", Handler: oh.MaxPlayersCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderprice", Handler: oh.PriceCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderjoin", Handler: oh.JoinCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderleave", Handler: oh.LeaveCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercancel", Handler: oh.CancelCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercancelcomfirm", Handler: oh.CancelComfirmCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordershow", Handler: oh.ShowCallback})
	}
	for _, handler := range oh.CallbackHandlers {
		result, err = handler.ProceedCallback(cq)
	}
	return
}

func (oh *OrderBotHandler) ProceedMessage(msg *telegram.Message) (result telegram.MessageResponse, err error) {
	if len(oh.MessageHandlers) == 0 {
		order_cmd := telegram.CommandHandler{
			Command: oh.OrderCommand, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
				return oh.CreateOrder(m, nil)
			}}
		oh.MessageHandlers = append(oh.MessageHandlers, &order_cmd)
		list_cmd := telegram.CommandHandler{
			Command: oh.ListCommand, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
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

	currTime := time.Now()
	stime := time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 8, 0, 0, 0, currTime.Location())
	etime := stime.Add(time.Duration(time.Hour))

	res, err := oh.OrderService.CreateOrder(reserve.Reserve{Person: p, StartTime: stime, EndTime: etime}, nil)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("creating order error: %s", err.Error()),
			AnswerMsg: "Can't create order"}
		return oh.SendNessageError(msg, err.(telegram.HelperError), nil)
	}

	var kbd telegram.InlineKeyboardMarkup
	kh := oh.GetReserveActions(res, *msg.From)
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

func (oh *OrderBotHandler) GetReserveActions(res reserve.Reserve, user telegram.User) (h telegram.KeyboardHelper) {
	ah := telegram.ActionsKeyboardHelper{Data: res.Id.String()}
	if res.Canceled {
		return &ah
	}
	ah.Columns = 2
	if res.Orderd() {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderjoin", Text: oh.Resources.JoinPlayer.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderleave", Text: oh.Resources.JoinPlayer.LeaveButton})
	}
	if res.Person.TelegramId == user.Id {
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderdate", Text: oh.Resources.DateTime.DateButton})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordertime", Text: oh.Resources.DateTime.TimeButton})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordersets", Text: oh.Resources.Set.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderminlevel", Text: oh.Resources.Level.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordercourts", Text: oh.Resources.Court.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderplayers", Text: oh.Resources.MaxPlayer.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "orderprice", Text: oh.Resources.Price.Button})
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordercancel", Text: oh.Resources.Cancel.Button})
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

	if len(rlist) == 0 {
		cq.Answer(oh.Bot, oh.Resources.NoReservesAnswer, nil)
		result = oh.Bot.SendMessage(&telegram.MessageRequest{
			Text:   oh.Resources.NoReservesMessage,
			ChatId: cq.Message.Chat.Id})
		return

	}
	for _, res := range rlist {
		kh := oh.GetReserveActions(res, *cq.From)
		kh.SetData(res.Id.String())
		mr := oh.GetReserveMR(res, kh)
		mr.ChatId = cq.Message.Chat.Id
		result = oh.Bot.SendMessage(&mr)
	}
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

		return oh.UpdateReserveCQ(res, cq)
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
		return oh.UpdateReserveCQ(res, cq)
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
		return oh.UpdateReserveCQ(res, cq)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) ShowCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.OrderActionsHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	kbd := oh.GetReserveActions(res, *cq.From)
	mr := oh.GetReserveEditMR(res, kbd)
	mr.ChatId = cq.Message.Chat.Id
	cq.Message.EditText(oh.Bot, "", &mr)
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
		return oh.UpdateReserveCQ(res, cq)
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
		return oh.UpdateReserveCQ(res, cq)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) MaxPlayersCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	ch := oh.PlayerCountHelper
	err = ch.Parse(cq.Data)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err := oh.GetDataReserve(ch.Data, nil)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}
	if ch.Action == "set" {
		res.MaxPlayers = ch.Count
		return oh.UpdateReserveCQ(res, cq)
	} else {
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
		return oh.UpdateReserveCQ(res, cq)
	} else {
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) JoinCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
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
		pl_count := 0
		for pl_id, pl := range res.Players {
			if pl_id != p.Id {
				pl_count += pl.Count
			}
		}
		ch.Min = 1
		ch.Max = res.MaxPlayers - pl_count

		mr := oh.GetReserveEditMR(res, &ch)
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
	res.Players[p.Id] = reserve.Player{Person: p, Count: count}
	return oh.UpdateReserveCQ(res, cq)
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
			mr := oh.GetReserveMR(res, nil)
			mr.ChatId = pl.Person.TelegramId
			oh.Bot.SendMessage(&mr)
		}
	}
	return oh.UpdateReserveCQ(res, cq)
}

func (oh *OrderBotHandler) UpdateReserveCQ(res reserve.Reserve, cq *telegram.CallbackQuery) (telegram.MessageResponse, error) {

	res, err := oh.UpdateReserve(res)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, *cq.From))
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
