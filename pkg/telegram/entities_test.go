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

func TestMessageCreateEditTextRequest(t *testing.T) {
	msg := Message{Chat: &Chat{Id: 123456789}, MessageId: 987}
	mer := msg.CreateEditTextRequest("Hello", nil)

	t.Run("Error is nil", func(t *testing.T) {
		if mer.ChatId != msg.Chat.Id {
			t.Fail()
		}
		if mer.MessageId != msg.MessageId {
			t.Fail()
		}
		if mer.Text != "Hello" {
			t.Fail()
		}
	})
}

func TestMessageCreateMessageRequest(t *testing.T) {
	msg := Message{Chat: &Chat{Id: 123456789}, MessageId: 987}
	mer := msg.CreateMessageRequest("Hello", nil)

	t.Run("Error is nil", func(t *testing.T) {
		if mer.ChatId != msg.Chat.Id {
			t.Fail()
		}
		if mer.Text != "Hello" {
			t.Fail()
		}
	})
}

func TestMessageGetCommand(t *testing.T) {
	tests := map[string]struct {
		text string
		want string
	}{
		"Text without command": {
			text: "some text",
			want: "",
		},
		"Command": {
			text: "/start",
			want: "start",
		},
		"Command with bot name": {
			text: "/start@bot",
			want: "start",
		},
		"Command with text": {
			text: "/start some text",
			want: "start",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := Message{Text: test.text}
			if msg.GetCommand() != test.want {
				t.Fail()
			}
		})
	}
}

func TestMessageIsCommand(t *testing.T) {
	tests := map[string]struct {
		text string
		want bool
	}{
		"Text without command": {
			text: "some text",
			want: false,
		},
		"Command": {
			text: "/start",
			want: true,
		},
		"Command with bot name": {
			text: "/start@bot",
			want: true,
		},
		"Command with text": {
			text: "/start some text",
			want: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := Message{Text: test.text}
			if msg.IsCommand() != test.want {
				t.Fail()
			}
		})
	}
}
