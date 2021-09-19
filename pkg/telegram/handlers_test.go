package telegram

import (
	"testing"
)

func TestUpdateHandlerProceed(t *testing.T) {
	tb := Bot{Token: "***Token***"}
	message := Message{}
	tu := Update{Message: &message}

	uh := BaseUpdateHandler{}
	mh := BaseMessageHandler{
		Handler: func(_ *Bot, tm *Message) (bool, error) {
			tm.Caption = "Test caption"
			return true, nil
		},
	}
	uh.MessageHandlers = append(uh.MessageHandlers, &mh)

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

func TestCommandHandlerProceed(t *testing.T) {
	tests := map[string]struct {
		text        string
		want        string
		want_regexp string
	}{
		"Text without command": {
			text:        "some text",
			want:        "",
			want_regexp: "",
		},
		"Command start": {
			text:        "/start",
			want:        "Ok",
			want_regexp: "Ok",
		},
		"Command start with bot name": {
			text:        "/start@bot",
			want:        "Ok",
			want_regexp: "Ok",
		},
		"Command with text": {
			text:        "/start some text",
			want:        "Ok",
			want_regexp: "Ok",
		},
		"Command starting with text": {
			text:        "/starting some text",
			want:        "",
			want_regexp: "Ok",
		},
		"Another command with text": {
			text:        "/stop some text",
			want:        "",
			want_regexp: "",
		},
	}

	mh := BaseMessageHandler{
		Handler: func(_ *Bot, tm *Message) (bool, error) {
			tm.Caption = "Ok"
			return true, nil
		},
	}

	ch := CommandHandler{Command: "start", InnerHandler: &mh}

	for name, test := range tests {
		t.Run("Simple "+name, func(t *testing.T) {
			msg := Message{Text: test.text, Caption: ""}
			ch.Proceed(nil, &msg)
			if msg.Caption != test.want {
				t.Fail()
			}
		})
	}

	ch.IsRegexp = true

	for name, test := range tests {
		t.Run("Regexp "+name, func(t *testing.T) {
			msg := Message{Text: test.text, Caption: ""}
			ch.Proceed(nil, &msg)
			if msg.Caption != test.want_regexp {
				t.Fail()
			}
		})
	}
}
