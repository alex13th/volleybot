package telegram

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseResponse(t *testing.T) {
	updateCount := 3
	json := `{
		"ok": true,
		"result": [
			{
				"update_id": 123130161,
				"message": {
					"message_id": 2468,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630134810,
					"text": "Hello world!!!"
				}
			},
			{
				"update_id": 123130162,
				"message": {
					"message_id": 2469,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630135377,
					"text": "\u041f\u0440\u0438\u0432\u0435\u0442!"
				}
			},
			{
				"update_id": 123130163,
				"message": {
					"message_id": 2470,
					"from": {"id": 586350636,"is_bot": false,"first_name": "Alexey","last_name": "Sukharev","language_code": "en"},
					"chat": {"id": 586350636,"first_name": "Alexey","last_name": "Sukharev","type": "private"},
					"date": 1630135517,
					"text": "Hello!"
				}
			}
		]
	}`

	expUpdate := Update{
		UpdateId: 123130161,
		Message: Message{
			MessageId: 2468,
			From: &User{
				Id:           586350636,
				IsBot:        false,
				FirstName:    "Alexey",
				LastName:     "Sukharev",
				LanguageCode: "en",
			},
			Chat: &Chat{
				Id:        586350636,
				FirstName: "Alexey",
				LastName:  "Sukharev",
				Type:      "private",
			},
			Date: 1630134810,
			Text: "Hello world!!!",
		},
	}

	resp := UpdateResponse{}
	err := resp.Parse(strings.NewReader(json))

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

	t.Run(fmt.Sprintf("Update count is %d", updateCount), func(t *testing.T) {
		if len(resp.Result) != 3 {
			t.Fail()
		}
	})

	t.Run("First update properties", func(t *testing.T) {
		if !reflect.DeepEqual(expUpdate, resp.Result[0]) {
			t.Fail()
		}
	})
}
