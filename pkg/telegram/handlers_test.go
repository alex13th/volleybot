package telegram

import (
	"testing"
)

func TestUpdateHandlerProceed(t *testing.T) {
	tb, _ := NewSimpleBot("***Token***", nil)
	message := Message{}
	tu := Update{Message: &message}

	uh := BaseUpdateHandler{}
	mh := BaseMessageHandler{
		Handler: func(tm *Message) error {
			tm.Caption = "Test caption"
			return nil
		},
	}
	uh.MessageHandlers = append(uh.MessageHandlers, &mh)

	uh.ProceedUpdate(tb, tu)

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
		Handler: func(tm *Message) error {
			tm.Caption = "Ok"
			return nil
		},
	}

	ch := CommandHandler{Command: "start", Handler: func(m *Message) error {
		return mh.ProceedMessage(m)
	}}

	for name, test := range tests {
		t.Run("Simple "+name, func(t *testing.T) {
			msg := Message{Text: test.text, Caption: ""}
			ch.ProceedMessage(&msg)
			if msg.Caption != test.want {
				t.Fail()
			}
		})
	}

	ch.IsRegexp = true

	for name, test := range tests {
		t.Run("Regexp "+name, func(t *testing.T) {
			msg := Message{Text: test.text, Caption: ""}
			ch.ProceedMessage(&msg)
			if msg.Caption != test.want_regexp {
				t.Fail()
			}
		})
	}
}

func TestPrefixCallbackHandlerProceed(t *testing.T) {
	tests := map[string]struct {
		text        string
		want        string
		want_regexp string
	}{
		"Data without prefix": {
			text: "some text",
			want: "some text",
		},
		"Prefix pref": {
			text: "pref_some",
			want: "Ok",
		},
		"Another prefixt": {
			text: "alternate_pref",
			want: "alternate_pref",
		},
	}

	handler := PrefixCallbackHandler{
		Prefix: "pref",
		Handler: func(cb *CallbackQuery) error {
			cb.Data = "Ok"
			return nil
		},
	}

	for name, test := range tests {
		t.Run("Callback "+name, func(t *testing.T) {
			callback := CallbackQuery{Data: test.text}
			handler.ProceedCallback(&callback)
			if callback.Data != test.want {
				t.Fail()
			}
		})
	}
}
