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
	}

	admin, _ := person.NewPerson("Admin")
	admin.LocationRoles[loc.Id] = []string{"admin"}
	admin.TelegramId = 321
	akbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.NewReserveBtn, CallbackData: "res_main_order"},
		},
	}

	tests := map[string]struct {
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"User chat": {p: person.Person{TelegramId: 100}, cid: 100, kbd: kbd},
		"Admin":     {p: admin, cid: admin.TelegramId, kbd: append(akbd, kbd...)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_main_main")
			bp, _ := NewBaseStateProvider(st, msg, test.p, loc, nil, "")
			sp := MainStateProvider{BaseStateProvider: bp, Resources: res}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}

func TestShowKbd(t *testing.T) {
	res := NewShowResourcesRu()
	plid, _ := uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	oauthor := person.Person{Id: plid, Firstname: "Elly", TelegramId: 100}
	plid, _ = uuid.Parse("80155587-168c-4255-82ec-991119f3e110")
	pl2 := person.Person{Id: plid, Firstname: "Steve", TelegramId: 123}
	plid, _ = uuid.Parse("da10db9a-490b-4010-9d8c-561cca979dd0")
	pl3 := person.Person{Id: plid, Firstname: "Tina", TelegramId: 123456}
	r := volley.Volley{Reserve: reserve.Reserve{
		Id:        uuid.New(),
		Location:  location.Location{Id: uuid.New()},
		Person:    oauthor,
		StartTime: time.Date(2021, 12, 04, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2021, 12, 04, 17, 0, 0, 0, time.UTC)},
		CourtCount: 1,
		MaxPlayers: 4,
		Players:    []person.Player{{Person: oauthor, Count: 2}, {Person: pl2, Count: 3}, {Person: pl3, Count: 1}},
	}
	admin, _ := person.NewPerson("Admin")
	admin.LocationRoles[r.Location.Id] = []string{"admin"}
	admin.TelegramId = 321

	tests := map[string]struct {
		res volley.Volley
		p   person.Person
		cid int
		kbd [][]telegram.InlineKeyboardButton
	}{
		"Group chat": {res: r, p: person.Person{}, cid: -10,
			kbd: [][]telegram.InlineKeyboardButton{{
				{Text: res.JoinBtn, CallbackData: "res_show_join_" + r.Id.String()},
				{Text: res.JoinLeaveBtn, CallbackData: "res_show_leave_" + r.Id.String()}},
				{{Text: res.RefreshBtn, CallbackData: "res_show_refresh_" + r.Id.String()}},
			}},
		"Not joined Person": {res: r, p: person.Person{TelegramId: 10}, cid: 10,
			kbd: [][]telegram.InlineKeyboardButton{{
				{Text: res.JoinBtn, CallbackData: "res_show_join_" + r.Id.String()},
				{Text: res.JoinMultiBtn, CallbackData: "res_show_joinm_" + r.Id.String()}},
				{{Text: res.RefreshBtn, CallbackData: "res_show_refresh_" + r.Id.String()}},
			}},
		"Joined Person": {res: r, p: person.Person{TelegramId: 123}, cid: 123,
			kbd: [][]telegram.InlineKeyboardButton{
				{
					{Text: res.JoinMultiBtn, CallbackData: "res_show_joinm_" + r.Id.String()},
					{Text: res.JoinTimeBtn, CallbackData: "res_show_jtime_" + r.Id.String()},
				},
				{
					{Text: res.JoinLeaveBtn, CallbackData: "res_show_leave_" + r.Id.String()},
					{Text: res.RefreshBtn, CallbackData: "res_show_refresh_" + r.Id.String()},
				},
			}},
		"Admin Person": {res: r, p: admin, cid: admin.TelegramId,
			kbd: [][]telegram.InlineKeyboardButton{
				{
					{Text: res.DateTime.DateBtn, CallbackData: "res_show_date_" + r.Id.String()},
					{Text: res.DateTime.TimeBtn, CallbackData: "res_show_time_" + r.Id.String()},
				},
				{
					{Text: res.SetsBtn, CallbackData: "res_show_sets_" + r.Id.String()},
					{Text: res.DescriptionBtn, CallbackData: "res_show_desc_" + r.Id.String()},
				},
				{
					{Text: res.SettingsBtn, CallbackData: "res_show_settings_" + r.Id.String()},
					{Text: res.ActionsBtn, CallbackData: "res_show_actions_" + r.Id.String()},
				},
				{
					{Text: res.JoinBtn, CallbackData: "res_show_join_" + r.Id.String()},
					{Text: res.JoinMultiBtn, CallbackData: "res_show_joinm_" + r.Id.String()},
				},
				{
					{Text: res.RefreshBtn, CallbackData: "res_show_refresh_" + r.Id.String()},
				},
			}},
		"Author Person Joined": {res: r, p: oauthor, cid: oauthor.TelegramId,
			kbd: [][]telegram.InlineKeyboardButton{
				{
					{Text: res.DateTime.DateBtn, CallbackData: "res_show_date_" + r.Id.String()},
					{Text: res.DateTime.TimeBtn, CallbackData: "res_show_time_" + r.Id.String()},
				},
				{
					{Text: res.SetsBtn, CallbackData: "res_show_sets_" + r.Id.String()},
					{Text: res.DescriptionBtn, CallbackData: "res_show_desc_" + r.Id.String()},
				},
				{
					{Text: res.SettingsBtn, CallbackData: "res_show_settings_" + r.Id.String()},
					{Text: res.ActionsBtn, CallbackData: "res_show_actions_" + r.Id.String()},
				},
				{
					{Text: res.JoinMultiBtn, CallbackData: "res_show_joinm_" + r.Id.String()},
					{Text: res.JoinTimeBtn, CallbackData: "res_show_jtime_" + r.Id.String()},
				},
				{
					{Text: res.JoinLeaveBtn, CallbackData: "res_show_leave_" + r.Id.String()},
					{Text: res.RefreshBtn, CallbackData: "res_show_refresh_" + r.Id.String()},
				},
			}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_show_show_" + test.res.Id.String())
			bp, _ := NewBaseStateProvider(st, msg, test.p, test.res.Location, nil, "")
			bp.reserve = test.res
			sp := ShowStateProvider{BaseStateProvider: bp, Resources: res}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}

func TestSettingsStateKbd(t *testing.T) {
	res := NewSettingsResourcesRu()
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
		Players:    []person.Player{{Person: oauthor, Count: 2}},
	}
	admin, _ := person.NewPerson("Admin")
	admin.LocationRoles[r.Location.Id] = []string{"admin"}
	admin.TelegramId = 321
	akbd := [][]telegram.InlineKeyboardButton{
		{
			{Text: res.ActivityBtn, CallbackData: "res_settings_activity_" + r.Id.String()},
			{Text: res.CourtBtn, CallbackData: "res_settings_courts_" + r.Id.String()},
		},
		{
			{Text: res.LevelBtn, CallbackData: "res_settings_level_" + r.Id.String()},
			{Text: res.MaxBtn, CallbackData: "res_settings_max_" + r.Id.String()},
		},
		{
			{Text: res.PriceBtn, CallbackData: "res_settings_price_" + r.Id.String()},
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
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_settings_settings_" + test.res.Id.String())
			bp, _ := NewBaseStateProvider(st, msg, test.p, test.res.Location, nil, "")
			bp.reserve = test.res
			sp := SettingsStateProvider{BaseStateProvider: bp, Resources: res}
			acts := sp.GetKeyboardHelper().GetKeyboard()
			if !reflect.DeepEqual(acts, test.kbd) {
				t.Fail()
			}
		})
	}
}

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
		Players:    []person.Player{{Person: oauthor, Count: 2}},
	}
	admin, _ := person.NewPerson("Admin")
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
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_actions_actions_" + test.res.Id.String())
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

func TestJoinmStateKbd(t *testing.T) {
	res := NewJoinPlayersResourcesRu()
	plid, _ := uuid.Parse("14a959fe-b3bb-4538-b7eb-feabc2f5c2c8")
	oauthor := person.Person{Id: plid, Firstname: "Elly", TelegramId: 100}
	r := volley.Volley{Reserve: reserve.Reserve{
		Person: oauthor},
		CourtCount: 1,
		MaxPlayers: 4,
		Players:    []person.Player{{Person: oauthor, Count: 2}},
	}

	tests := map[string]struct {
		res   volley.Volley
		p     person.Person
		cid   int
		count int
	}{
		"Group chat":        {res: r, p: person.Person{}, cid: -10, count: 2},
		"Group chat joined": {res: r, p: oauthor, cid: -10, count: 2},
		"Chat":              {res: r, p: person.Person{}, cid: oauthor.TelegramId, count: 4},
		"Chat joined":       {res: r, p: oauthor, cid: oauthor.TelegramId, count: 4},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_joinm_joinm_" + test.res.Id.String())
			bp, _ := NewBaseStateProvider(st, msg, test.p, test.res.Location, nil, "")
			bp.reserve = test.res
			sp := JoinPlayersStateProvider{BaseStateProvider: bp, Resources: res}
			kbd := sp.GetKeyboardHelper().(*telegram.CountKeyboardHelper)

			if kbd.Max != test.count {
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
		Players:    []person.Player{{Person: oauthor, Count: 2}},
	}
	admin, _ := person.NewPerson("Admin")
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
			msg := telegram.Message{Chat: &telegram.Chat{Id: test.cid}}
			st, _ := telegram.NewState().Parse("res_cancel_cancel_" + test.res.Id.String())
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
