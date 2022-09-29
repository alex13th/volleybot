package bvbot

import (
	"reflect"
	"testing"
	"time"

	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

func TestActionsStateKbd(t *testing.T) {
	res := NewActionsResourcesRu()
	plid, _ := uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	oauthor := person.Person{Id: plid, Firstname: "Elly", TelegramId: 100}
	r := volley.Volley{Reserve: reserve.Reserve{
		Id:        uuid.New(),
		Location:  location.Location{Id: uuid.New()},
		Person:    oauthor,
		StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC)},
		CourtCount: 1,
		MaxPlayers: 4,
		Members:    []volley.Member{{Player: volley.Player{Person: oauthor}, Count: 2}},
	}
	admin := person.NewPerson("Admin")
	admin.LocationRoles[r.Location.Id] = []string{"admin"}
	admin.TelegramId = 321
	akbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.CancelBtn, CallbackData: "res_actions_cancel_" + r.Id.String()},
			{Text: res.CopyBtn, CallbackData: "res_actions_copy_" + r.Id.String()},
		},
		{
			{Text: res.PublishBtn, CallbackData: "res_actions_pub_" + r.Id.String()},
			{Text: res.RemovePlayerBtn, CallbackData: "res_actions_rmpl_" + r.Id.String()},
		},
	}

	tests := map[string]struct {
		res volley.Volley
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"Group chat":       {res: r, p: person.Person{}, cid: -10},
		"Group chat admin": {res: r, p: admin, cid: -10},
		"Admin":            {res: r, p: admin, cid: admin.TelegramId, kbd: akbd},
		"Author":           {res: r, p: oauthor, cid: oauthor.TelegramId, kbd: akbd},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{}
			st, _ := telegram.NewState().Parse("res_actions_actions_" + test.res.Id.String())
			st.ChatId = test.cid
			bp, _ := NewBaseStateProvider(st, msg, test.p, test.res.Location, nil, "")
			bp.reserve = test.res
			sp := ActionsStateProvider{BaseStateProvider: bp, Resources: res}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}

func TestCancelStateKbd(t *testing.T) {
	res := NewCancelResourcesRu()
	sres := NewShowResourcesRu()
	plid, _ := uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	oauthor := person.Person{Id: plid, Firstname: "Elly", TelegramId: 100}
	r := volley.Volley{Reserve: reserve.Reserve{
		Id:        uuid.New(),
		Location:  location.Location{Id: uuid.New()},
		Person:    oauthor,
		StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC)},
		CourtCount: 1,
		MaxPlayers: 4,
		Members:    []volley.Member{{Player: volley.Player{Person: oauthor}, Count: 2}},
	}
	admin := person.NewPerson("Admin")
	admin.LocationRoles[r.Location.Id] = []string{"admin"}
	admin.TelegramId = 321
	akbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.ConfirmBtn, CallbackData: "res_cancel_confirm_" + r.Id.String()},
			{Text: sres.JoinLeaveBtn, CallbackData: "res_cancel_leave_" + r.Id.String()},
		},
	}

	tests := map[string]struct {
		res volley.Volley
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"Group chat":       {res: r, p: person.Person{}, cid: -10},
		"Group chat admin": {res: r, p: admin, cid: -10},
		"Admin":            {res: r, p: admin, cid: admin.TelegramId, kbd: akbd},
		"Author":           {res: r, p: oauthor, cid: oauthor.TelegramId, kbd: akbd},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{}
			st, _ := telegram.NewState().Parse("res_cancel_cancel_" + test.res.Id.String())
			st.ChatId = test.cid
			bp, _ := NewBaseStateProvider(st, msg, test.p, test.res.Location, nil, "")
			bp.reserve = test.res
			sp := CancelStateProvider{BaseStateProvider: bp, Resources: res, ShowResources: sres}

			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}
