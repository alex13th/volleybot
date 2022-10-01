package bvbot

import (
	"reflect"
	"testing"

	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

func TestMainStateKbd(t *testing.T) {
	res := NewMainResourcesRu()
	loc := location.Location{Id: uuid.New()}

	kbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.TodayBtn, CallbackData: "res_main_today"},
		},
		{
			{Text: res.ListDateBtn, CallbackData: "res_main_listd"},
		},
		{
			{Text: res.ProfileBtn, CallbackData: "res_main_profile"},
		},
	}

	admin := person.NewPerson("Admin")
	admin.LocationRoles[loc.Id] = []string{"admin"}
	admin.TelegramId = 321
	okbd := [][]telegram.InlineKeyboardButton{
		{{Text: res.NewReserveBtn, CallbackData: "res_main_order"}},
	}

	akbd := []telegram.InlineKeyboardButton{{Text: res.ConfigBtn, CallbackData: "res_main_config"}}

	tests := map[string]struct {
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"User chat": {p: person.Person{TelegramId: 100}, cid: 100, kbd: kbd},
		"Admin":     {p: admin, cid: admin.TelegramId, kbd: append(append(okbd, kbd...), akbd)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{}
			st, _ := telegram.NewState().Parse("res_main_main")
			st.ChatId = test.cid
			bp, _ := NewBaseStateProvider(st, msg, test.p, loc, nil, nil, "")
			sp := MainStateProvider{BaseStateProvider: bp, Resources: res}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}
