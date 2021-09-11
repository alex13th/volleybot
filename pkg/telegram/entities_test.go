package telegram

import (
	"testing"
)

func TestMessageCreateReplyRequest(t *testing.T) {
	msg := Message{Chat: &Chat{Id: 123456789}, MessageId: 987}
	mr := msg.CreateReplyRequest("Hello", nil)

	t.Run("Error is nil", func(t *testing.T) {
		if mr.ChatId != msg.Chat.Id {
			t.Fail()
		}
		if mr.ReplyToMessageId != msg.MessageId {
			t.Fail()
		}
		if mr.Text != "Hello" {
			t.Fail()
		}
	})
}
