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

func TestGetUpdatesParams(t *testing.T) {
	tests := map[string]struct {
		request UpdatesRequest
		want    string
	}{
		"Empty parameters": {
			request: UpdatesRequest{},
			want:    "",
		},
		"Fully filled parameters": {
			request: UpdatesRequest{
				Offset:         551,
				Limit:          100,
				Timeout:        20,
				AllowedUpdates: []string{"message", "edited_channel_post", "callback_query"},
			},
			want: "allowed_updates=message&allowed_updates=edited_channel_post&allowed_updates=callback_query&limit=100&offset=551&timeout=20",
		},
	}

	req := tests["Fully filled parameters"].request
	_, method, err := req.GetParams()

	t.Run("Error is nil", func(t *testing.T) {
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Method", func(t *testing.T) {
		if method != "getUpdates" {
			t.Fail()
		}
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req := test.request
			values, _, err := req.GetParams()

			if err != nil {
				t.Fail()
			}

			if values.Encode() != test.want {
				t.Fail()
			}
		})
	}

}
