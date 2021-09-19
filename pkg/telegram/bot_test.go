package telegram

import (
	"io"
	"testing"
)

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

func TestBotSendRequest(t *testing.T) {
	tb, _ := NewBot(Bot{Token: "***Token***"})
	tb.Client = httpClientMock{}

	req := MessageRequest{
		ChatId: 586350636,
		Text:   "Message text",
	}
	resp, err := tb.SendRequest(&req)

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

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

func TestBotSendMessage(t *testing.T) {
	tb, _ := NewBot(Bot{Token: "***Token***"})
	tb.Client = httpClientMock{
		Body: `{
			"ok": true,
			"result": {
				"message_id": 2468,
				"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
				"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
				"date": 1630134810,
				"text": "Hello world!!!"
			}
		}`,
	}

	req := &MessageRequest{
		ChatId: 586350636,
		Text:   "Message text",
	}

	botResp, err := tb.SendMessage(req)

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Response Ok", func(t *testing.T) {
		if !botResp.Ok {
			t.Fail()
		}
	})

	t.Run("Response Message Id", func(t *testing.T) {
		if botResp.Result.MessageId != 2468 {
			t.Fail()
		}
	})
}

func TestBotGetUpdates(t *testing.T) {
	tb := Bot{Token: "***Token***"}
	tb.Client = httpClientMock{Body: `{
		"ok": true,
		"result": [{"update_id": 123130161},{"update_id": 123130162},{"update_id": 123130163}]
	}`}

	botResp, err := tb.GetUpdates()

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Response Ok", func(t *testing.T) {
		if !botResp.Ok {
			t.Fail()
		}
	})

	t.Run("New offset", func(t *testing.T) {
		if tb.Request.Offset != 123130164 {
			t.Fail()
		}
	})
}

func TestUpdateHandlerProceed(t *testing.T) {
	tb := Bot{Token: "***Token***"}
	message := Message{}
	tu := Update{Message: &message}

	uh := UpdateHandler{}
	mh := MessageHandler{
		Handler: func(_ *Bot, tm *Message) error {
			tm.Caption = "Test caption"
			return nil
		},
	}
	uh.MessageHandlers = append(uh.MessageHandlers, mh)

	err := uh.Proceed(&tb, tu)

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Message proceeded", func(t *testing.T) {
		if message.Caption != "Test caption" {
			t.Fail()
		}
	})

}
