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

type Bot struct {
	ApiEndpoint string
	Token       string
	Client      HttpClient
	Parameters  UpdatesRequest
}

func (tb *Bot) GetUpdates() (botResp UpdateResponse, err error) {
	httpResp, err := tb.SendRequest(&tb.Parameters)

	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	err = botResp.Parse(httpResp.Body)

	if len(botResp.Result) > 0 {
		tb.Parameters.Offset = botResp.Result[len(botResp.Result)-1].UpdateId + 1
	}

	return
}

func NewBot(Settings Bot) (*Bot, error) {
	if Settings.ApiEndpoint == "" {
		Settings.ApiEndpoint = DefaultApiUrl
	}
	return &Bot{ApiEndpoint: Settings.ApiEndpoint, Token: Settings.Token}, nil
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

func (tb *Bot) SendMessage(msg MessageRequest) (botResp MessageResponse, err error) {
	httpResp, err := tb.SendRequest(&msg)

	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	botResp = MessageResponse{}
	err = botResp.Parse(httpResp.Body)
	return
}
