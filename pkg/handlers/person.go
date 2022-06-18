package handlers

import (
	"fmt"
	"log"
	"strconv"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"
)

func NewPersonHandler(tb *telegram.Bot, serv *services.PersonService, rl PersonResourceLoader) (h PersonBotHandler) {
	h = PersonBotHandler{Bot: tb, PersonService: *serv}
	h.Resources = rl.GetResource()
	levels := []telegram.EnumItem{}
	for i := 0; i <= 80; i += 10 {
		levels = append(levels, telegram.EnumItem{Id: strconv.Itoa(i),
			Item: fmt.Sprintf("%s %s", person.PlayerLevel(i).Emoji(), person.PlayerLevel(i))})
	}
	h.LevelHelper = telegram.NewEnumKeyboardHelper(h.Resources.Level.Message, "personlevel", levels)
	sexs := []telegram.EnumItem{
		{Id: "1", Item: fmt.Sprintf("%s %s", person.Sex(1).Emoji(), person.Sex(1))},
		{Id: "2", Item: fmt.Sprintf("%s %s", person.Sex(2).Emoji(), person.Sex(2))},
	}
	h.SexHelper = telegram.NewEnumKeyboardHelper(h.Resources.Level.Message, "personsex", sexs)
	return
}

type PersonResources struct {
	ProfileCommand telegram.BotCommand
	Level          PlayerLevelResources
}

type PersonResourceLoader struct{}

func (rl PersonResourceLoader) GetResource() (r PersonResources) {
	r.ProfileCommand.Command = "profile"
	r.ProfileCommand.Description = "настройки профиля пользователя"
	return
}

type PersonBotHandler struct {
	Bot              *telegram.Bot
	PersonService    services.PersonService
	Resources        PersonResources
	LevelHelper      telegram.EnumKeyboardHelper
	SexHelper        telegram.EnumKeyboardHelper
	StateRepository  telegram.StateRepository
	MessageHandlers  []telegram.MessageHandler
	CallbackHandlers []telegram.CallbackHandler
}

func (h *PersonBotHandler) GetCommands(tuser *telegram.User) (cmds []telegram.BotCommand) {
	cmds = append(cmds, h.Resources.ProfileCommand)
	return
}

func (h *PersonBotHandler) ProceedMessage(msg *telegram.Message) (result telegram.MessageResponse, err error) {
	if msg.Chat.Id <= 0 {
		return
	}
	if len(h.MessageHandlers) == 0 {
		profile_cmd := telegram.CommandHandler{
			Command: h.Resources.ProfileCommand.Command, Handler: func(m *telegram.Message) (telegram.MessageResponse, error) {
				return h.ShowProfile(m, nil)
			},
		}
		h.MessageHandlers = append(h.MessageHandlers, &profile_cmd)
		for _, handler := range h.MessageHandlers {
			result, err = handler.ProceedMessage(msg)
		}
		return
	}
	for _, handler := range h.MessageHandlers {
		result, err = handler.ProceedMessage(msg)
	}
	return
}

func (h *PersonBotHandler) ProceedCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	if len(h.CallbackHandlers) == 0 {
		h.CallbackHandlers = append(h.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "personlevel", Handler: h.LevelCallback})
		h.CallbackHandlers = append(h.CallbackHandlers,
			&telegram.PrefixCallbackHandler{Prefix: "personsex", Handler: h.SexCallback})
	}
	for _, handler := range h.CallbackHandlers {
		result, err = handler.ProceedCallback(cq)
	}
	return
}

func (h *PersonBotHandler) ShowProfile(msg *telegram.Message, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	p, _ := h.PersonService.GetByTelegramId(msg.From.Id)
	mr := h.GetPersonMR(p, h.GetPersonActions(p))

	return msg.SendMessage(h.Bot, "", &mr), nil
}

func (h *PersonBotHandler) LevelCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	p, _ := h.GetPerson(cq.From)
	ch := h.LevelHelper
	ch.Parse(cq.Data)

	if ch.Action == "set" {
		lvl, _ := strconv.Atoi(ch.Choice)
		p.Level = person.PlayerLevel(lvl)
		return h.UpdatePersonCQ(p, cq, false)
	} else {
		mr := h.GetPersonEditMR(p, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(h.Bot, "", &mr)
		return cq.Answer(h.Bot, "Ok", nil), nil
	}
}

func (h *PersonBotHandler) SexCallback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	p, _ := h.GetPerson(cq.From)
	ch := h.SexHelper
	ch.Parse(cq.Data)

	if ch.Action == "set" {
		si, _ := strconv.Atoi(ch.Choice)
		p.Sex = person.Sex(si)
		return h.UpdatePersonCQ(p, cq, false)
	} else {
		mr := h.GetPersonEditMR(p, &ch)
		mr.ChatId = cq.Message.Chat.Id
		cq.Message.EditText(h.Bot, "", &mr)
		return cq.Answer(h.Bot, "Ok", nil), nil
	}
}

func (h *PersonBotHandler) GetPerson(tuser *telegram.User) (p person.Person, err error) {
	p, err = h.PersonService.GetByTelegramId(tuser.Id)
	if err != nil {
		log.Println(err.Error())
		_, ok := err.(person.ErrorPersonNotFound)
		if ok {
			p, _ = person.NewPerson(tuser.FirstName)
			p.TelegramId = tuser.Id
			p.Lastname = tuser.LastName
			p, err = h.PersonService.Add(p)
		}
	}
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("getting person error: %s", err.Error()),
			AnswerMsg: "Can't get person"}
	}
	return
}

func (h *PersonBotHandler) UpdatePersonCQ(p person.Person, cq *telegram.CallbackQuery, renew bool) (resp telegram.MessageResponse, err error) {

	h.UpdatePerson(p)

	mr := h.GetPersonEditMR(p, h.GetPersonActions(p))
	cq.Message.EditText(h.Bot, "", &mr)
	return cq.Answer(h.Bot, "Ok", nil), nil
}

func (h *PersonBotHandler) UpdatePerson(p person.Person) error {
	return h.PersonService.Update(p)
}

func (oh *PersonBotHandler) GetPersonEditMR(p person.Person, kh telegram.KeyboardHelper) (mer telegram.EditMessageTextRequest) {
	mr := oh.GetPersonMR(p, kh)
	return telegram.EditMessageTextRequest{Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (h *PersonBotHandler) GetPersonMR(p person.Person, kh telegram.KeyboardHelper) (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	if kh != nil {
		kh.SetData(p.Id.String())
		kbd.InlineKeyboard = append(kbd.InlineKeyboard, kh.GetKeyboard()...)
		kbdText = "\n*" + kh.GetText() + "* "
	}

	pview := person.NewTelegramViewRu(p)
	mtxt := fmt.Sprintf("%s\n%s", pview.GetText(), kbdText)
	if len(kbd.InlineKeyboard) > 0 {
		return telegram.MessageRequest{Text: mtxt, ParseMode: pview.ParseMode, ReplyMarkup: kbd}
	}
	return telegram.MessageRequest{Text: mtxt, ParseMode: pview.ParseMode}
}

func (h *PersonBotHandler) GetPersonActions(p person.Person) (лh telegram.KeyboardHelper) {
	ah := telegram.ActionsKeyboardHelper{Data: p.Id.String()}
	ah.Columns = 2

	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "personlevel", Text: "Уровень"})
	ah.Actions = append(ah.Actions, telegram.ActionButton{
		Prefix: "personsex", Text: "Пол"})
	return &ah
}
