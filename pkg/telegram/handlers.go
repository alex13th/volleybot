package telegram

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
