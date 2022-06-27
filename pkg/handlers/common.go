package handlers

import (
	"fmt"
	"log"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"
)

type BotHandler interface {
	GetCallbackHandlers() []telegram.CallbackHandler
	GetMessageHandler() []telegram.MessageHandler
	GetCommands(tuser *telegram.User) (cmds []telegram.BotCommand)
}

type CommonHandler struct {
	Bot             *telegram.Bot
	StateRepository telegram.StateRepository
	PersonService   *services.PersonService
}

func (h *CommonHandler) SendCallbackError(cq *telegram.CallbackQuery, cq_err telegram.HelperError, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	log.Println(cq_err.Error())
	result = cq.Answer(h.Bot, cq_err.AnswerMsg, nil)
	if chanr != nil {
		chanr <- result
	}
	return result, cq_err
}

func (h *CommonHandler) SendMessageError(msg *telegram.Message, m_err telegram.HelperError, chanr chan telegram.MessageResponse) (result telegram.MessageResponse, err error) {
	log.Println(m_err.Error())
	result = msg.Reply(h.Bot, m_err.AnswerMsg, nil)
	if chanr != nil {
		chanr <- result
	}
	return result, m_err
}

func (h *CommonHandler) GetPersonCq(cq *telegram.CallbackQuery) (p person.Person, resp telegram.MessageResponse, err error) {
	p, err = h.GetPerson(cq.From)
	if err != nil {
		resp, err = h.SendCallbackError(cq, err.(telegram.HelperError), nil)
		return
	}
	return
}

func (h *CommonHandler) GetPerson(tuser *telegram.User) (p person.Person, err error) {
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
