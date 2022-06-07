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

type DescriptionResources struct {
	Message string
	Button  string
}

type JoinPlayerResources struct {
	Message     string
	Button      string
	LeaveButton string
	MultiButton string
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
	Location          location.Location
	BackBtn           string
	RefreshBtn        string
	PublishBtn        string
	ListCommand       telegram.BotCommand
	OrderCommand      telegram.BotCommand
	Locale            monday.Locale
	Description       DescriptionResources
	DateTime          DateTimeResources
	Court             CourtResources
	Level             PlayerLevelResources
	Set               SetResources
	MaxPlayer         MaxPlayerResources
	JoinPlayer        JoinPlayerResources
	Price             PriceResources
	Cancel            CancelResources
	ReservesMessage   string
	NoReservesMessage string
	NoReservesAnswer  string
	OkAnswer          string
}

type OrderResourceLoader interface {
	GetResource() OrderResources
}

type DefaultResourceLoader struct{}

func (rl DefaultResourceLoader) GetResource() (or OrderResources) {
	or.ListCommand.Command = "list"
	or.ListCommand.Description = "поиск подходящей площадки"
	or.OrderCommand.Command = "order"
	or.OrderCommand.Description = "заказать площадку(и)"
	or.BackBtn = "Назад"
	or.RefreshBtn = "Обновить"
	or.PublishBtn = "Опубликовать"
	or.Description.Button = "Описание"
	or.Description.Message = "❓Каким будет описание❓"
	or.Locale = monday.LocaleRuRU
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
	or.Court.Max = 6
	or.Court.MaxPlayers = 6
	or.MaxPlayer.Message = "❓Максимальное количество игроков❓"
	or.MaxPlayer.Button = "👫 Мест"
	or.MaxPlayer.Min = 4
	or.MaxPlayer.Max = or.Court.Max * or.Court.MaxPlayers
	or.JoinPlayer.Message = "❓Сколько игроков записать❓"
	or.JoinPlayer.Button = "😀 Буду"
	or.JoinPlayer.MultiButton = "🤩 Буду не один"
	or.JoinPlayer.LeaveButton = "😞 Не смогу"
	or.Price.Message = "❓Почем будет поиграть❓"
	or.Price.Button = "💰 Стоимость"
	or.Price.Min = 0
	or.Price.Max = 2000
	or.Price.Step = 100
	or.Cancel.Button = "💥Отменить"
	or.Cancel.Message = fmt.Sprintf("\n🧨*ВНИМАНИЕ!!!*🧨\nИгра будет отменена для всех участников. Если есть желание только выписаться, лучше воспользоваться кнопкой \"%s\"",
		or.JoinPlayer.LeaveButton)
	or.Cancel.Confirm = "🧨 Уверен"
	or.Cancel.Abort = "👌 Передумал"
	or.ReservesMessage = "❓Какую запись показать ❓"
	or.NoReservesMessage = "На дату %s нет доступных записей."
	or.NoReservesAnswer = "Резервы отсутствуют"
	or.OkAnswer = "Ok"

	return
}

func NewOrderHandler(tb *telegram.Bot, os *services.OrderService, rl OrderResourceLoader) (oh OrderBotHandler) {
	oh = OrderBotHandler{Bot: tb, OrderService: os}
	oh.Resources = rl.GetResource()

	oh.DateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderdate")
	oh.ListDateHelper = telegram.NewDateKeyboardHelper(oh.Resources.DateTime.DateMessage, "orderlistdate")
	oh.TimeHelper = telegram.NewTimeKeyboardHelper(oh.Resources.DateTime.TimeMessage, "ordertime")

	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i), Item: reserve.PlayerLevel(i)})
	}
	oh.MinLevelHelper = telegram.NewEnumKeyboardHelper(oh.Resources.Level.Message, "orderminlevel", levels)

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
	StateRepository    telegram.StateRepository
	Resources          OrderResources
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

func (oh *OrderBotHandler) GetCommands() (cmds []telegram.BotCommand) {
	cmds = append(cmds, oh.Resources.ListCommand)
	return
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
			&telegram.PrefixCallbackHandler{Prefix: "orderjoinmult", Handler: oh.JoinMultiCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderleave", Handler: oh.LeaveCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercancel", Handler: oh.CancelCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordercancelcomfirm", Handler: oh.CancelComfirmCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "ordershow", Handler: oh.ShowCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderlist", Handler: oh.ListOrdersCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderpub", Handler: oh.PublishCallback})
		oh.CallbackHandlers = append(oh.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "orderdesc", Handler: oh.DescriptionCallback})

	}
	for _, handler := range oh.CallbackHandlers {
		result, err = handler.ProceedCallback(cq)
	}
	return
}

func (oh *OrderBotHandler) ProceedMessage(msg *telegram.Message) (result telegram.MessageResponse, err error) {
	if msg.Chat.Id <= 0 {
		return
	}
	if len(oh.MessageHandlers) == 0 {
		order_cmd := telegram.CommandHandler{
			Command: oh.Resources.OrderCommand.Command, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
				return oh.CreateOrder(m, nil)
			}}
		oh.MessageHandlers = append(oh.MessageHandlers, &order_cmd)
		list_cmd := telegram.CommandHandler{
			Command: oh.Resources.ListCommand.Command, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
				return oh.ListOrders(m, nil), nil
			}}
		oh.MessageHandlers = append(oh.MessageHandlers, &list_cmd)
		players_state := telegram.StateMessageHandler{State: "orderplayers", StateRepository: oh.StateRepository,
			Handler: oh.MaxPlayersState}
		oh.MessageHandlers = append(oh.MessageHandlers, &players_state)
		desc_state := telegram.StateMessageHandler{State: "orderdesc", StateRepository: oh.StateRepository,
			Handler: oh.DescriptionState}
		oh.MessageHandlers = append(oh.MessageHandlers, &desc_state)
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

func (oh *OrderBotHandler) SendMessageError(msg *telegram.Message, m_err telegram.HelperError, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
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
	kh := oh.GetReserveActions(res, p, msg.Chat.Id)
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

func (oh *OrderBotHandler) GetReserveActions(res reserve.Reserve, p person.Person, chid int) (h telegram.KeyboardHelper) {
	ah := telegram.ActionsKeyboardHelper{Data: res.Id.String()}
	if res.Canceled {
		return &ah
	}
	ah.Columns = 2
	if res.Orderd() {
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
	}
	if chid == p.TelegramId {
		if res.Person.TelegramId == p.TelegramId || p.CheckLocationRole(res.Location, "admin") {
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
				Prefix: "orderdesc", Text: oh.Resources.Description.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "ordercancel", Text: oh.Resources.Cancel.Button})
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderpub", Text: oh.Resources.PublishBtn})
		}
	}
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "ordershow", Text: oh.Resources.RefreshBtn})

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
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
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
	kbd := oh.GetReserveActions(res, p, cq.Message.Chat.Id)
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
	kbd := oh.GetReserveActions(res, person.Person{}, cq.Message.Chat.Id)
	mr := oh.GetReserveMR(res, kbd)
	mr.ChatId = res.Location.ChatId
	oh.Bot.SendMessage(&mr)
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
		ch.Max = res.Location.CourtCount
		mr := oh.GetReserveEditMR(res, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(oh.Bot, "", &mr)
		return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
	}
}

func (oh *OrderBotHandler) DescriptionState(msg *telegram.Message, state telegram.State) (result telegram.MessageResponse, err error) {
	res, err := oh.GetDataReserve(state.Data, nil)
	if err != nil {
		oh.StateRepository.Clear(state.ChatId)
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	res.Description = msg.Text
	oh.StateRepository.Clear(state.ChatId)
	return oh.UpdateReserveMsg(res, msg, state.MessageId)
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
	mr := oh.GetReserveEditMR(res, &ah)
	mr.ChatId = cq.Message.Chat.Id
	cq.Message.EditText(oh.Bot, "", &mr)
	return cq.Answer(oh.Bot, oh.Resources.OkAnswer, nil), nil
}

func (oh *OrderBotHandler) MaxPlayersState(msg *telegram.Message, state telegram.State) (result telegram.MessageResponse, err error) {
	count, err := strconv.Atoi(msg.Text)
	if err != nil {
		herr := telegram.HelperError{Msg: err.Error(), AnswerMsg: "Maximum players count message convert error"}
		oh.StateRepository.Clear(state.ChatId)
		return oh.SendMessageError(msg, herr, nil)
	}

	res, err := oh.GetDataReserve(state.Data, nil)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	res.CourtCount = count
	oh.StateRepository.Clear(state.ChatId)
	return oh.UpdateReserveMsg(res, msg, state.MessageId)
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
		oh.StateRepository.Clear(cq.Message.Chat.Id)
		return oh.UpdateReserveCQ(res, cq)
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
		return oh.UpdateReserveCQ(res, cq)
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
		pl_count := 0
		for pl_id, pl := range res.Players {
			if pl_id != p.Id {
				pl_count += pl.Count
			}
		}
		ch.Min = 1
		ch.Max = res.MaxPlayers - pl_count
		var mr telegram.EditMessageTextRequest
		if ch.Max > 1 {
			mr = oh.GetReserveEditMR(res, &ch)
		} else {
			oh.GetReserveActions(res, p, cq.Message.Chat.Id)
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

func (oh *OrderBotHandler) NotifyPlayers(res reserve.Reserve, id int) {
	for _, pl := range res.Players {
		if pl.Person.TelegramId != id {
			p, _ := oh.OrderService.Persons.GetByTelegramId(pl.Person.TelegramId)
			if param, ok := p.Settings["notify"]; ok && param == "on" {
				mr := oh.GetReserveMR(res, nil)
				mr.ChatId = pl.Person.TelegramId
				oh.Bot.SendMessage(&mr)
			}
		}
	}
}

func (oh *OrderBotHandler) UpdateReserveCQ(res reserve.Reserve, cq *telegram.CallbackQuery) (resp telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(cq.From)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	res, err = oh.UpdateReserve(res)
	if err != nil {
		return oh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, p, cq.Message.Chat.Id))
	cq.Message.EditText(oh.Bot, "", &mr)
	oh.NotifyPlayers(res, cq.From.Id)

	return cq.Answer(oh.Bot, "Ok", nil), nil
}

func (oh *OrderBotHandler) UpdateReserveMsg(res reserve.Reserve, msg *telegram.Message, mid int) (resp telegram.MessageResponse, err error) {
	p, err := oh.GetPerson(msg.From)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	res, err = oh.UpdateReserve(res)
	if err != nil {
		return oh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}

	mr := oh.GetReserveEditMR(res, oh.GetReserveActions(res, p, msg.Chat.Id))
	if mid > 0 {
		mr.ChatId = msg.Chat.Id
		mr.MessageId = mid
	}
	resp = oh.Bot.SendMessage(&mr)
	oh.NotifyPlayers(res, msg.From.Id)

	return
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
