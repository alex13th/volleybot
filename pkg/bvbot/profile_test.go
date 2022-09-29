package bvbot

import (
	"reflect"
	"testing"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

func TestProfileStateKbd(t *testing.T) {
	res := NewProfileResourcesRu()
	loc := location.Location{Id: uuid.New()}

	kbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.LevelBtn, CallbackData: "res_profile_plevel"},
			{Text: res.SexBtn, CallbackData: "res_profile_sex"},
		},
		{
			{Text: res.NotifiesBtn, CallbackData: "res_profile_notifies"},
		},
	}

	admin := person.NewPerson("Admin")
	admin.LocationRoles[loc.Id] = []string{"admin"}
	admin.TelegramId = 321

	tests := map[string]struct {
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"User chat": {p: person.Person{TelegramId: 100}, cid: 100, kbd: kbd},
		"Admin":     {p: admin, cid: admin.TelegramId, kbd: kbd},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{}
			st, _ := telegram.NewState().Parse("res_profile_profile")
			st.ChatId = test.cid
			bp, _ := NewBaseStateProvider(st, msg, test.p, loc, nil, "")
			pp := PlayerStateProvider{BaseStateProvider: bp, Resources: res}
			sp := ProfileStateProvider{PlayerStateProvider: pp}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}
