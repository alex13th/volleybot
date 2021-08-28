package telegram

import (
	"testing"
)

func TestGetMessageParams(t *testing.T) {
	botMessage := MessageRequest{
		ChatId: 586350636,
		Text:   "Example of text",
	}

	values, method, err := botMessage.GetParams()

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Method", func(t *testing.T) {
		if method != "sendMessage" {
			t.Fail()
		}
	})

	t.Run("Required parameters", func(t *testing.T) {
		if values["chat_id"][0] != "586350636" || values["text"][0] != "Example of text" {
			t.Fail()
		}
	})
}
