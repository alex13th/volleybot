package telegram

import "regexp"

type UpdateHandler interface {
	Proceed(tb *Bot, update Update) error
	AppendMessageHandler(mh MessageHandler)
}

type BaseUpdateHandler struct {
	MessageHandlers []MessageHandler
}

func (handler *BaseUpdateHandler) AppendMessageHandler(mh MessageHandler) {
	handler.MessageHandlers = append(handler.MessageHandlers, mh)
}

func (uh BaseUpdateHandler) Proceed(tb *Bot, update Update) error {
	if update.Message != nil {
		for _, handler := range uh.MessageHandlers {
			cont, err := handler.Proceed(tb, update.Message)
			if !cont || err != nil {
				return err
			}
		}
	}
	return nil
}

type MessageHandler interface {
	Proceed(tb *Bot, tm *Message) (bool, error)
}

type BaseMessageHandler struct {
	Handler func(*Bot, *Message) (bool, error)
}

func (mh *BaseMessageHandler) Proceed(tb *Bot, tm *Message) (bool, error) {
	return mh.Handler(tb, tm)
}

type CommandHandler struct {
	InnerHandler MessageHandler
	Command      string
	IsRegexp     bool
}

func (mh *CommandHandler) Proceed(tb *Bot, tm *Message) (bool, error) {
	if mh.IsRegexp {
		re, err := regexp.Compile(mh.Command)
		if err != nil {
			return true, err
		}
		cmd := tm.GetCommand()
		if re.MatchString(cmd) {
			return mh.InnerHandler.Proceed(tb, tm)
		}
	}
	if tm.GetCommand() == mh.Command {
		return mh.InnerHandler.Proceed(tb, tm)
	}
	return true, nil
}
