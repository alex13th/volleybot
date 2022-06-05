package telegram

import (
	"errors"
	"regexp"
	"strings"
)

type CallbackQueryFunc func(cq *CallbackQuery) (MessageResponse, error)
type MessageFunc func(m *Message) (MessageResponse, error)
type MessageStateFunc func(m *Message, state State) (MessageResponse, error)

type UpdateHandler interface {
	ProceedUpdate(tb *Bot, update Update)
	AppendMessageHandler(mh MessageHandler)
	AppendCallbackHandler(ch CallbackHandler)
}

type BaseUpdateHandler struct {
	MessageHandlers  []MessageHandler
	CallbackHandlers []CallbackHandler
}

func (handler *BaseUpdateHandler) AppendCallbackHandler(ch CallbackHandler) {
	handler.CallbackHandlers = append(handler.CallbackHandlers, ch)
}

func (handler *BaseUpdateHandler) AppendMessageHandler(mh MessageHandler) {
	handler.MessageHandlers = append(handler.MessageHandlers, mh)
}

func (uh BaseUpdateHandler) ProceedUpdate(tb *Bot, update Update) {
	if update.Message != nil {
		for _, handler := range uh.MessageHandlers {
			handler.ProceedMessage(update.Message)
		}
	}
	if update.CallbackQuery != nil {
		for _, handler := range uh.CallbackHandlers {
			handler.ProceedCallback(update.CallbackQuery)
		}
	}
}

type CallbackHandler interface {
	ProceedCallback(*CallbackQuery) (MessageResponse, error)
}

type BaseCallbackHandler struct {
	Bot     *Bot
	Handler CallbackQueryFunc
}

func (h *BaseCallbackHandler) ProceedCallback(cb *CallbackQuery) (MessageResponse, error) {
	return h.Handler(cb)
}

type PrefixCallbackHandler struct {
	Bot     *Bot
	Prefix  string
	Handler CallbackQueryFunc
}

func (h *PrefixCallbackHandler) ProceedCallback(cb *CallbackQuery) (MessageResponse, error) {
	var prefix []string = strings.Split(cb.Data, "_")
	if len(prefix) > 1 && prefix[0] == h.Prefix {
		return h.Handler(cb)
	}
	return MessageResponse{}, errors.New("data hasn't prefix")
}

type MessageHandler interface {
	ProceedMessage(tm *Message) (MessageResponse, error)
}

type BaseMessageHandler struct {
	Bot     *Bot
	Handler MessageFunc
}

func (h *BaseMessageHandler) ProceedMessage(m *Message) (MessageResponse, error) {
	return h.Handler(m)
}

type CommandHandler struct {
	Bot      *Bot
	Handler  MessageFunc
	Command  string
	Commands []BotCommand
	IsRegexp bool
}

func (h *CommandHandler) GetCommands() []BotCommand {
	return h.Commands
}

func (h *CommandHandler) ProceedMessage(m *Message) (result MessageResponse, err error) {
	if h.IsRegexp {
		var re *regexp.Regexp
		re, err = regexp.Compile(h.Command)
		if err != nil {
			return
		}
		cmd := m.GetCommand()
		if re.MatchString(cmd) {
			return h.Handler(m)
		}
	}
	if m.GetCommand() == h.Command {
		return h.Handler(m)
	}
	return
}

type StateMessageHandler struct {
	Bot             *Bot
	Handler         MessageStateFunc
	State           string
	StateRepository StateRepository
}

func (h *StateMessageHandler) ProceedMessage(m *Message) (result MessageResponse, err error) {
	state, err := h.StateRepository.Get(m.Chat.Id)
	if state.State == h.State {
		return h.Handler(m, state)
	}
	return
}
