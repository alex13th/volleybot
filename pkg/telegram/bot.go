package telegram

import (
	"encoding/json"
	"fmt"
	"io"
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

type Response struct {
	Ok          bool               `json:"ok"`
	Result      []Update           `json:"result"`
	Description string             `json:"description"`
	ErrorCode   int                `json:"error_code"`
	Parameters  ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
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

func ParseResponse(reader io.Reader) (response Response, err error) {
	dec := json.NewDecoder(reader)
	err = dec.Decode(&response)
	return
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

func (tb Bot) SendMessage(msg MessageRequest) error {
	_, err := tb.SendRequest(&msg)
	return err
}
