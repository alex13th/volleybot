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
}

type GetUpdateParams struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type LongPoller struct {
	Parameters GetUpdateParams
	Bot        Bot
	Client     HttpClient
}

type BotMessage struct {
	Bot     Bot
	Message Message
}

func NewBot(Settings Bot) (*Bot, error) {
	if Settings.ApiEndpoint == "" {
		Settings.ApiEndpoint = DefaultApiUrl
	}
	return &Bot{ApiEndpoint: Settings.ApiEndpoint, Token: Settings.Token}, nil
}

func (tb *Bot) SendRequest(request Request) (httpResp *http.Response, err error) {
	values, method, _ := request.GetParams()
	url := fmt.Sprintf("%s/bot%s/%s", tb.ApiEndpoint, tb.Token, method)

	httpReq, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err = tb.Client.Do(httpReq)
	return
}

func (tb Bot) SendMessage(msg MessageRequest) (response MessageResponse, err error) {
	httpResp, err := tb.SendRequest(&msg)

	if err != nil {
		return
	}
	response = MessageResponse{}
	err = response.Parse(httpResp.Body)
	return
}
