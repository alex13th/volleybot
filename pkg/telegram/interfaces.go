package telegram

import (
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Logger interface {
	Printf(format string, v ...any)
}

type Bot interface {
	GetUpdates(UpdatesRequest) (*UpdateResponse, error)
	SendRequest(Request) (*http.Response, error)
	SendMessage(Request) (*MessageResponse, error)
}

type LongPoller interface {
	Run() error
}

type Poller interface {
	ProceedUpdates(resp *UpdateResponse)
}

type UpdateHandler interface {
	ProceedUpdate(tb Bot, update Update) error
	AppendMessageHandlers(...MessageHandler)
	AppendCallbackHandlers(...CallbackHandler)
}

type CallbackHandler interface {
	ProceedCallback(*CallbackQuery) error
}

type MessageHandler interface {
	ProceedMessage(tm *Message) error
}

type MessageRequestHelper interface {
	GetEditMR() EditMessageTextRequest
	GetMR() MessageRequest
}

type CallbackDataParser interface {
	GetAction() string
	GetPrefix() string
	GetState() State
	GetValue() string
	Parse(string) error
	SetState(state State)
}

type KeyboardHelper interface {
	GetKeyboard() interface{}
	GetText() string
}

type Request interface {
	GetParams() (url.Values, string, error)
}

type Response interface {
	Parse(reader io.Reader) error
}

type StateBuilder interface {
	GetStateProvider(State) (StateProvider, error)
}

type StateProvider interface {
	GetRequests() []StateRequest
	Proceed() (State, error)
}

type StateRepository interface {
	Get(ChatId int) ([]State, error)
	GetByData(Data string) ([]State, error)
	GetByMessage(msg Message) (State, error)
	Set(State) error
	Clear(State) error
}
