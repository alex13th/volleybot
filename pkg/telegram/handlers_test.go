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
