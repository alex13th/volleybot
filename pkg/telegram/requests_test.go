package telegram

import (
	"testing"
)

func TestGetMessageParams(t *testing.T) {
	tests := map[string]struct {
		request MessageRequest
		want    map[string]string
	}{
		"Required parameters": {
			request: MessageRequest{
				ChatId: 586350636,
				Text:   "Example of text",
			},
			want: map[string]string{
				"chat_id": "586350636",
				"text":    "Example of text",
			},
		},
		"Fully filled parameters": {
			request: MessageRequest{
				ChatId:    586350636,
				Text:      "Example of text",
				ParseMode: "MarkdownV2",
				Entities: []MessageEntity{
					{
						Type:   "url",
						Offset: 0,
						Length: 5,
						Url:    "https://google.com",
					},
					{
						Type:   "mention",
						Offset: 6,
						Length: 5,
						User: &User{
							Id:        987654321,
							IsBot:     false,
							FirstName: "Firstname",
						},
					},
				},
				DisableWebPagePreview:    true,
				DisableNotification:      true,
				ReplyToMessageId:         1234,
				AllowSendingWithoutReply: true,
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: []InlineKeyboardButton{
						{Text: "Button text 1", CallbackData: "Data1"},
						{Text: "Button text 2", CallbackData: "Data2"},
					},
				},
			},
			want: map[string]string{
				"allow_sending_without_reply": "true",
				"chat_id":                     "586350636",
				"disable_notification":        "true",
				"disable_web_page_preview":    "true",
				"entities":                    `[{"type":"url","offset":0,"length":5,"url":"https://google.com","user":null,"language":""},{"type":"mention","offset":6,"length":5,"url":"","user":{"id":987654321,"is_bot":false,"first_name":"Firstname","last_name":"","username":"","language_code":"","can_join_groups":false,"can_read_all_group_messages":false,"supports_inline_queries":false},"language":""}]`,
				"parse_mode":                  "MarkdownV2",
				"reply_to_message_id":         "1234",
				"text":                        "Example of text",
				"reply_markup":                `{"inline_keyboard":[{"text":"Button text 1","url":"","callback_data":"Data1","switch_inline_query":"","switch_inline_query_current_chat":"","pay":false},{"text":"Button text 2","url":"","callback_data":"Data2","switch_inline_query":"","switch_inline_query_current_chat":"","pay":false}]}`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()

			if err != nil {
				t.Fail()
			}

			if method != "sendMessage" {
				t.Fail()
			}

			for name, val := range test.want {
				valStr := values.Get(name)
				if valStr != val {
					t.Fail()
				}
			}
		})
	}
}

func TestGetEditMessageTextParams(t *testing.T) {
	tests := map[string]struct {
		request EditMessageTextRequest
		want    map[string]string
	}{
		"Required parameters": {
			request: EditMessageTextRequest{
				ChatId:    586350636,
				MessageId: 123456789,
				Text:      "Example of text",
			},
			want: map[string]string{
				"chat_id":    "586350636",
				"message_id": "123456789",
				"text":       "Example of text",
			},
		},
		"Fully filled parameters": {
			request: EditMessageTextRequest{
				ChatId:    586350636,
				MessageId: 123456789,
				Text:      "Example of text",
				ParseMode: "MarkdownV2",
				Entities: []MessageEntity{
					{
						Type:   "url",
						Offset: 0,
						Length: 5,
						Url:    "https://google.com",
					},
					{
						Type:   "mention",
						Offset: 6,
						Length: 5,
						User: &User{
							Id:        987654321,
							IsBot:     false,
							FirstName: "Firstname",
						},
					},
				},
				DisableWebPagePreview: true,
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: []InlineKeyboardButton{
						{Text: "Button text 1", CallbackData: "Data1"},
						{Text: "Button text 2", CallbackData: "Data2"},
					},
				},
			},
			want: map[string]string{
				"chat_id":                  "586350636",
				"message_id":               "123456789",
				"disable_web_page_preview": "true",
				"entities":                 `[{"type":"url","offset":0,"length":5,"url":"https://google.com","user":null,"language":""},{"type":"mention","offset":6,"length":5,"url":"","user":{"id":987654321,"is_bot":false,"first_name":"Firstname","last_name":"","username":"","language_code":"","can_join_groups":false,"can_read_all_group_messages":false,"supports_inline_queries":false},"language":""}]`,
				"parse_mode":               "MarkdownV2",
				"text":                     "Example of text",
				"reply_markup":             `{"inline_keyboard":[{"text":"Button text 1","url":"","callback_data":"Data1","switch_inline_query":"","switch_inline_query_current_chat":"","pay":false},{"text":"Button text 2","url":"","callback_data":"Data2","switch_inline_query":"","switch_inline_query_current_chat":"","pay":false}]}`,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			values, method, err := test.request.GetParams()

			if err != nil {
				t.Fail()
			}

			if method != "editMessageText" {
				t.Fail()
			}

			for name, val := range test.want {
				valStr := values.Get(name)
				if valStr != val {
					t.Fail()
				}
			}
		})
	}
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
