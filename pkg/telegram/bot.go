package telegram

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var DefaultApiUrl string = "https://api.telegram.org"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewBot(Settings *Bot) (*Bot, error) {
	if Settings.ApiEndpoint == "" {
		Settings.ApiEndpoint = DefaultApiUrl
	}
	return &Bot{
		ApiEndpoint: Settings.ApiEndpoint,
		Token:       Settings.Token,
		ChatStates:  make(map[int]interface{}),
	}, nil
}

type Bot struct {
	ApiEndpoint string
	Token       string
	Client      HttpClient
	Request     UpdatesRequest
	ChatStates  map[int]interface{}
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
	if request == nil {
		return httpResp, errors.New("request in nil")
	}
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

func (tb *Bot) SendMessage(req Request) (resp MessageResponse, err error) {
	httpResp, err := tb.SendRequest(req)

	if err == nil {
		defer httpResp.Body.Close()
		err = resp.Parse(httpResp.Body)
	}
	if err != nil {
		resp.Description = err.Error()
		return resp, err
	}
	return resp, err
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
				handler.ProceedUpdate(lp.Bot, update)
			}
		}
	}
}
