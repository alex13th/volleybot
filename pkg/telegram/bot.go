package telegram

import (
	"fmt"
	"net/http"
	"strings"
)

var DefaultApiUrl string = "https://api.telegram.org"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewBot(Settings Bot) (*Bot, error) {
	if Settings.ApiEndpoint == "" {
		Settings.ApiEndpoint = DefaultApiUrl
	}
	return &Bot{ApiEndpoint: Settings.ApiEndpoint, Token: Settings.Token}, nil
}

type Bot struct {
	ApiEndpoint string
	Token       string
	Client      HttpClient
	Request     UpdatesRequest
}

func (tb *Bot) GetUpdates() (botResp UpdateResponse, err error) {
	httpResp, err := tb.SendRequest(&tb.Request)

	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	err = botResp.Parse(httpResp.Body)

	if len(botResp.Result) > 0 {
		tb.Request.Offset = botResp.Result[len(botResp.Result)-1].UpdateId + 1
	}

	return
}

func (tb *Bot) SendRequest(request Request) (httpResp *http.Response, err error) {
	values, method, err := request.GetParams()

	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/bot%s/%s", tb.ApiEndpoint, tb.Token, method)

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err = tb.Client.Do(httpReq)
	return
}

func (tb *Bot) SendMessage(msg Request) (botResp MessageResponse, err error) {
	httpResp, err := tb.SendRequest(msg)

	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	botResp = MessageResponse{}
	err = botResp.Parse(httpResp.Body)
	return
}

func (tb *Bot) NewPoller() (lp LongPoller, err error) {
	uh := BaseUpdateHandler{MessageHandlers: []MessageHandler{}}
	lp = LongPoller{Bot: tb, UpdateHandlers: []UpdateHandler{&uh}}
	return
}

type LongPoller struct {
	Bot            *Bot
	UpdateHandlers []UpdateHandler
}

func (lp *LongPoller) Run() {
	for {
		updates, _ := lp.Bot.GetUpdates()
		for _, update := range updates.Result {
			for _, handler := range lp.UpdateHandlers {
				handler.Proceed(lp.Bot, update)
			}

		}
	}
}

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
