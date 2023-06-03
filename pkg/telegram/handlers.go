package telegram

import (
	"errors"
	"regexp"
	"strings"
)

type CallbackQueryFunc func(cq *CallbackQuery) error
type MessageFunc func(m *Message) error
type MessageStateFunc func(m *Message, state State) error

type BaseUpdateHandler struct {
	MessageHandlers  []MessageHandler
	CallbackHandlers []CallbackHandler
}

func (handler *BaseUpdateHandler) AppendCallbackHandlers(ch ...CallbackHandler) {
	handler.CallbackHandlers = append(handler.CallbackHandlers, ch...)
}

func (handler *BaseUpdateHandler) AppendMessageHandlers(mh ...MessageHandler) {
	handler.MessageHandlers = append(handler.MessageHandlers, mh...)
}

func (uh BaseUpdateHandler) ProceedUpdate(tb Bot, update Update) (err error) {
	if update.Message != nil {
		for _, handler := range uh.MessageHandlers {
			if err = handler.ProceedMessage(update.Message); err != nil {
				return
			}
		}
	}
	if update.CallbackQuery != nil {
		for _, handler := range uh.CallbackHandlers {
			if err = handler.ProceedCallback(update.CallbackQuery); err != nil {
				return
			}
		}
	}
	return
}

type BaseCallbackHandler struct {
	Bot     Bot
	Handler CallbackQueryFunc
}

func (h *BaseCallbackHandler) ProceedCallback(cb *CallbackQuery) error {
	return h.Handler(cb)
}

type PrefixCallbackHandler struct {
	Bot     Bot
	Prefix  string
	Handler CallbackQueryFunc
}

func (h *PrefixCallbackHandler) ProceedCallback(cb *CallbackQuery) error {
	var prefix []string = strings.Split(cb.Data, "_")
	if len(prefix) > 1 && prefix[0] == h.Prefix {
		return h.Handler(cb)
	}
	return errors.New("data hasn't prefix")
}

type BaseMessageHandler struct {
	Bot     Bot
	Handler MessageFunc
}

func (h *BaseMessageHandler) ProceedMessage(m *Message) error {
	return h.Handler(m)
}

type CommandHandler struct {
	Bot      Bot
	Handler  MessageFunc
	Command  string
	Commands []BotCommand
	IsRegexp bool
}

func (h *CommandHandler) GetCommands() []BotCommand {
	return h.Commands
}

func (h *CommandHandler) ProceedMessage(m *Message) (err error) {
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
	Bot             Bot
	Handler         MessageStateFunc
	State           string
	StateRepository StateRepository
}

func (h *StateMessageHandler) ProceedMessage(m *Message) error {
	slist, err := h.StateRepository.Get(m.Chat.Id)
	if err != nil {
		return err
	}
	for _, st := range slist {
		if st.State == h.State {
			return h.Handler(m, st)
		}
	}
	return err
}
