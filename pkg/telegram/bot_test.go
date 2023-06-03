package telegram

import (
	"errors"
	"io"
	"testing"
)

func TestSimpleBotSendRequest(t *testing.T) {
	tb, _ := NewSimpleBot("***Token***", httpClientMock{})
	tb.client = httpClientMock{}

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

func TestSimpleBotSendMessage(t *testing.T) {
	tb, _ := NewSimpleBot("***Token***",
		httpClientMock{
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
		})

	req := MessageRequest{
		ChatId: 586350636,
		Text:   "Message text",
	}

	botResp, _ := tb.SendMessage(req)

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

func TestSimpleBotGetUpdates(t *testing.T) {
	tb, _ := NewSimpleBot("***Token***",
		httpClientMock{Body: `{
		"ok": true,
		"result": [{"update_id": 123130161},{"update_id": 123130162},{"update_id": 123130163}]
	}`})

	resp, err := tb.GetUpdates(UpdatesRequest{})

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Response Ok", func(t *testing.T) {
		if !resp.Ok {
			t.Fail()
		}
	})

	t.Run("Response updates", func(t *testing.T) {
		if (resp.Result[0] != Update{UpdateId: 123130161} ||
			resp.Result[1] != Update{UpdateId: 123130162} ||
			resp.Result[2] != Update{UpdateId: 123130163}) {
			t.Fail()
		}
	})
}

func TestSimplePollerProceedUpdates(t *testing.T) {
	tb, _ := NewSimpleBot("***Token***",
		httpClientMock{Body: `{
		"ok": true,
		"result": [{"update_id": 123130161},{"update_id": 123130162},{"update_id": 123130163}]
	}`})

	pol := NewSimplePoller(tb)
	errHand := UpdateHandlerMock{err: errors.New("Mock error")}
	pol.UpdateHandlers = append(pol.UpdateHandlers, errHand)
	pol.Logger = LoggerMock{}
	err := pol.ProceedUpdates()

	t.Run("Check error", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Check offset", func(t *testing.T) {
		if pol.offset != 123130164 {
			t.Fail()
		}
	})
}
