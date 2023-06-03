package telegram

import (
	"fmt"
	"net/http"
	"strings"
)

const DefaultApiUrl string = "https://api.telegram.org"

func NewSimpleBot(Token string, client HttpClient) (SimpleBot, error) {
	return SimpleBot{
		apiEndpoint: DefaultApiUrl,
		client:      client,
		token:       Token,
		chatStates:  make(map[int]interface{}),
	}, nil
}

type SimpleBot struct {
	apiEndpoint string
	token       string
	client      HttpClient
	chatStates  map[int]interface{}
}

func (tb SimpleBot) GetUpdates(req UpdatesRequest) (resp *UpdateResponse, err error) {
	var httpResp *http.Response
	if httpResp, err = tb.SendRequest(req); err == nil {
		defer httpResp.Body.Close()
		resp = &UpdateResponse{}
		err = resp.Parse(httpResp.Body)
	}
	return
}

func (tb SimpleBot) SendRequest(botReq Request) (resp *http.Response, err error) {
	values, method, err := botReq.GetParams()
	if err == nil {
		var httpReq *http.Request
		url := fmt.Sprintf("%s/bot%s/%s", tb.apiEndpoint, tb.token, method)
		httpReq, err = http.NewRequest("POST", url, strings.NewReader(values.Encode()))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err = tb.client.Do(httpReq)
		}
	}
	return
}

func (tb SimpleBot) SendMessage(req Request) (resp *MessageResponse, err error) {
	var httpResp *http.Response
	if httpResp, err = tb.SendRequest(req); err == nil {
		defer httpResp.Body.Close()
		resp = &MessageResponse{}
		err = resp.Parse(httpResp.Body)
	}
	return
}

type SimplePoller struct {
	bot    Bot
	offset int
	Logger
	UpdateHandlers []UpdateHandler
}

func NewSimplePoller(tb Bot) SimplePoller {
	uh := BaseUpdateHandler{MessageHandlers: []MessageHandler{}}
	return SimplePoller{bot: tb, UpdateHandlers: []UpdateHandler{&uh}}
}

func (lp *SimplePoller) ProceedUpdates() (err error) {
	if resp, err := lp.bot.GetUpdates(UpdatesRequest{Offset: lp.offset}); err == nil {
		if len(resp.Result) > 0 {
			lp.offset = resp.Result[len(resp.Result)-1].UpdateId + 1
		}
		for _, update := range resp.Result {
			for _, handler := range lp.UpdateHandlers {
				if err = handler.ProceedUpdate(lp.bot, update); err != nil && lp.Logger != nil {
					lp.Logger.Printf("ERROR: SimplePoller proceed update error '%s'", err.Error())
				}
			}
		}
	}
	return nil
}

type SimpleLongPoller struct {
	SimplePoller
}

func (lp SimpleLongPoller) Run() error {
	for {
		lp.ProceedUpdates()
	}
}
