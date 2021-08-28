package telegram

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type httpClientMock struct{}

func (client httpClientMock) Do(httpRequest *http.Request) (*http.Response, error) {
	httpResponse := http.Response{}
	httpResponse.Request = httpRequest
	return &httpResponse, nil
}

func TestNewBot(t *testing.T) {
	tests := []struct {
		name     string
		settings Bot
		want     Bot
	}{
		{
			"Default Url test",
			Bot{Token: "SOME123TOKEN"},
			Bot{Token: "SOME123TOKEN", ApiEndpoint: "https://api.telegram.org"},
		},
		{
			"Custom Url test",
			Bot{Token: "SOME123TOKEN", ApiEndpoint: "new.api.telegram.com"},
			Bot{Token: "SOME123TOKEN", ApiEndpoint: "new.api.telegram.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb, err := NewBot(test.settings)
			if tb.ApiEndpoint != test.want.ApiEndpoint || tb.Token != test.settings.Token || err != nil {
				t.Fail()
			}

		})
	}
}

func TestParseResponse(t *testing.T) {
	updateCount := 3
	ResponseJson := `{
		"ok": true,
		"result": [
			{
				"update_id": 123130161,
				"message": {
					"message_id": 2468,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630134810,
					"text": "Hello world!!!"
				}
			},
			{
				"update_id": 123130162,
				"message": {
					"message_id": 2469,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630135377,
					"text": "\u041f\u0440\u0438\u0432\u0435\u0442!"
				}
			},
			{
				"update_id": 123130163,
				"message": {
					"message_id": 2470,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630135517,
					"text": "Hello!"
				}
			}
		]
	}`

	expUpdate := Update{
		UpdateId: 123130161,
		Message: Message{
			MessageId: 2468,
			From: &User{
				Id:           586350636,
				IsBot:        false,
				FirstName:    "Alexey",
				LastName:     "Sukharev",
				LanguageCode: "en",
			},
			Chat: &Chat{
				Id:        586350636,
				FirstName: "Alexey",
				LastName:  "Sukharev",
				Type:      "private",
			},
			Date: 1630134810,
			Text: "Hello world!!!",
		}}

	response, err := ParseResponse(strings.NewReader(ResponseJson))

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})
	t.Run("Response Ok", func(t *testing.T) {
		if !response.Ok {
			t.Fail()
		}
	})

	t.Run(fmt.Sprintf("Update count is %d", updateCount), func(t *testing.T) {
		if len(response.Result) != 3 {
			t.Fail()
		}
	})
	t.Run("First update properties", func(t *testing.T) {
		if !reflect.DeepEqual(expUpdate, response.Result[0]) {
			t.Fail()
		}
	})
}

func TestSendRequest(t *testing.T) {
	tb, _ := NewBot(Bot{Token: "***Token***"})
	tb.Client = httpClientMock{}

	req := MessageRequest{
		ChatId: 586350636,
		Text:   "Message text",
	}
	resp, _ := tb.SendRequest(&req)

	t.Run("Request method", func(t *testing.T) {
		if resp.Request.Method != "POST" {
			t.Fail()
		}
	})

	t.Run("Request URL", func(t *testing.T) {
		if resp.Request.URL.String() != "https://api.telegram.org/bot***Token***/sendMessage" {
			t.Fail()
		}
	})

	t.Run("Headers", func(t *testing.T) {
		if resp.Request.Header["Content-Type"][0] != "application/x-www-form-urlencoded" {
			t.Fail()
		}
	})

	t.Run("Request body", func(t *testing.T) {
		bytes, err := io.ReadAll(resp.Request.Body)
		if err != nil || string(bytes) != "chat_id=586350636&text=Message+text" {
			t.Fail()
		}
	})
}
