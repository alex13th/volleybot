package handlers

import (
	"fmt"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/res"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type ReserveMessager struct {
	Reserve        reserve.Reserve
	KeyboardHelper telegram.KeyboardHelper
	Resources      res.OrderResources
}

func NewReserveMessager(r reserve.Reserve, kh telegram.KeyboardHelper, res res.OrderResources) ReserveMessager {
	return ReserveMessager{Reserve: r, KeyboardHelper: kh, Resources: res}
}

func (rm *ReserveMessager) GetEditMR(ChatId int) (mer telegram.EditMessageTextRequest) {
	mr := rm.GetMR(ChatId)
	return telegram.EditMessageTextRequest{ChatId: ChatId, Text: mr.Text, ParseMode: mr.ParseMode, ReplyMarkup: mr.ReplyMarkup}
}

func (rm *ReserveMessager) GetMR(ChatId int) (mr telegram.MessageRequest) {
	var kbd telegram.InlineKeyboardMarkup
	var kbdText string
	if rm.KeyboardHelper != nil {
		rm.KeyboardHelper.SetData(rm.Reserve.Id.String())
		kbd.InlineKeyboard = append(kbd.InlineKeyboard, rm.KeyboardHelper.GetKeyboard()...)
		kbdText = "\n*" + rm.KeyboardHelper.GetText() + "* "
	}

	rview := reserve.NewTelegramViewRu(rm.Reserve)
	mtxt := fmt.Sprintf("%s\n%s", rview.GetText(), kbdText)
	if ChatId < 0 {
		mtxt += rm.Resources.MaxPlayer.GroupChatWarning
	}

	if len(kbd.InlineKeyboard) > 0 {
		return telegram.MessageRequest{ChatId: ChatId, Text: mtxt, ParseMode: rview.ParseMode, ReplyMarkup: kbd}
	}
	return telegram.MessageRequest{ChatId: ChatId, Text: mtxt, ParseMode: rview.ParseMode}
}

func (rm *ReserveMessager) SetReserveActions(p person.Person, ChatId int, state string) {
	ah := telegram.ActionsKeyboardHelper{Data: rm.Reserve.Id.String()}
	if rm.Reserve.Canceled {
		rm.KeyboardHelper = nil
	}
	ah.Columns = 2
	if ChatId == p.TelegramId {
		if rm.Reserve.Person.TelegramId == p.TelegramId || p.CheckLocationRole(rm.Reserve.Location, "admin") {
			if state == "ordershow" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderdate", Text: rm.Resources.DateTime.DateButton})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordertime", Text: rm.Resources.DateTime.TimeButton})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordersets", Text: rm.Resources.Set.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderdesc", Text: rm.Resources.Description.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordersettings", Text: rm.Resources.SettingsBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderactions", Text: rm.Resources.ActionsBtn})
			} else if state == "ordersettings" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderactivity", Text: rm.Resources.Activity.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercourts", Text: rm.Resources.Court.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderminlevel", Text: rm.Resources.Level.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderplayers", Text: rm.Resources.MaxPlayer.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderprice", Text: rm.Resources.Price.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordershow", Text: rm.Resources.BackBtn})
			} else if state == "orderactions" {
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercancel", Text: rm.Resources.Cancel.Button})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordercopy", Text: rm.Resources.CopyBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderpub", Text: rm.Resources.PublishBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "orderremovepl", Text: rm.Resources.RemovePlayerBtn})
				ah.Actions = append(ah.Actions, telegram.ActionButton{
					Prefix: "ordershow", Text: rm.Resources.BackBtn})
			}
		}
	}
	if rm.Reserve.Ordered() && state == "ordershow" {
		if ChatId <= 0 || !rm.Reserve.HasPlayerByTelegramId(p.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderjoin", Text: rm.Resources.JoinPlayer.Button})
		}
		if ChatId > 0 || rm.Reserve.MaxPlayers-rm.Reserve.PlayerCount(uuid.Nil) > 1 {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderjoinmult", Text: rm.Resources.JoinPlayer.MultiButton})
		}
		if ChatId > 0 && rm.Reserve.HasPlayerByTelegramId(p.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderarrivetime", Text: rm.Resources.JoinPlayer.ArriveButton})
		}
		if ChatId <= 0 || rm.Reserve.HasPlayerByTelegramId(p.TelegramId) {
			ah.Actions = append(ah.Actions, telegram.ActionButton{
				Prefix: "orderleave", Text: rm.Resources.JoinPlayer.LeaveButton})
		}
		ah.Actions = append(ah.Actions, telegram.ActionButton{
			Prefix: "ordershow", Text: rm.Resources.RefreshBtn})
	}
	rm.KeyboardHelper = &ah
}

type ReserveHandler struct {
	CommonHandler
	PlayerHandler *PlayerHandler
	Reserves      reserve.ReserveRepository
}

func (rh *ReserveHandler) GetDataReserve(data string, kh telegram.KeyboardHelper,
	rchan chan services.ReserveResult) (r reserve.Reserve, err error) {
	if kh != nil {
		if err = kh.Parse(data); err != nil {
			return
		}
		data = kh.GetData()
	}
	var id uuid.UUID
	id, err = uuid.Parse(data)
	if err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order getting reserve error: %s", err.Error()),
			AnswerMsg: "Parse reserve id error"}

	} else {
		r, err = rh.Reserves.Get(id)
		if err != nil {
			err = telegram.HelperError{
				Msg:       fmt.Sprintf("Getting reserve error: %s", err.Error()),
				AnswerMsg: "Getting reserve error"}
		}
	}

	if rchan != nil {
		rchan <- services.ReserveResult{Reserve: r, Err: err}
	}
	return
}

func (rh *ReserveHandler) UpdateReserveCQ(res reserve.Reserve, cq *telegram.CallbackQuery, state string, renew bool) (resp telegram.MessageResponse, err error) {
	st := telegram.State{
		ChatId:    cq.Message.Chat.Id,
		Data:      res.Id.String(),
		State:     state,
		MessageId: cq.Message.MessageId,
	}

	if err = rh.UpdateReserve(&res); err != nil {
		return rh.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rm := NewReserveMessager(res, nil, rh.Resources)
	rh.StateRepository.Set(st)
	rh.UpdateReserveMessages(res, true)
	rh.PlayerHandler.NotifyPlayers(res, &rm, cq.From.Id, "notify")
	resp = cq.Answer(rh.Bot, "Ok", nil)
	return resp, nil
}

func (rh *ReserveHandler) UpdateReserveMsg(res reserve.Reserve, msg *telegram.Message, mid int) (resp telegram.MessageResponse, err error) {
	st := telegram.State{Data: res.Id.String(), State: "ordershow"}
	if err = rh.UpdateReserve(&res); err != nil {
		return rh.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	rh.StateRepository.Set(st)
	rm := NewReserveMessager(res, nil, rh.Resources)
	rh.UpdateReserveMessages(res, true)
	rh.PlayerHandler.NotifyPlayers(res, &rm, msg.From.Id, "notify")
	return
}

func (rm *ReserveHandler) UpdateReserve(res *reserve.Reserve) (err error) {

	if err = rm.Reserves.Update(*res); err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order update can't update reserve %s error: %s", res.Id, err.Error()),
			AnswerMsg: "Can't update reserve"}
		return
	}
	*res, err = rm.Reserves.Get(res.Id)
	return
}

func (rh *ReserveHandler) UpdateReserveMessages(res reserve.Reserve, renew bool) {
	slist, _ := rh.StateRepository.GetByData(res.Id.String())
	notified := map[int]bool{}
	for _, st := range slist {
		if notified[st.ChatId] && st.ChatId < 0 {
			continue
		}
		notified[st.ChatId] = true
		p, _ := rh.GetPerson(&telegram.User{Id: st.ChatId})
		rm := NewReserveMessager(res, nil, rh.Resources)
		rm.SetReserveActions(p, st.ChatId, st.State)
		if renew && st.ChatId < 0 {
			mr := rm.GetMR(st.ChatId)
			mr.DisableNotification = true
			resp := rh.Bot.SendMessage(&mr)
			rh.Bot.SendMessage(&telegram.DeleteMessageRequest{ChatId: st.ChatId, MessageId: st.MessageId})
			rh.StateRepository.Clear(st)
			st.MessageId = resp.Result.MessageId
			rh.StateRepository.Set(st)
		} else {
			mr := rm.GetEditMR(st.ChatId)
			mr.MessageId = st.MessageId
			rh.Bot.SendMessage(&mr)
		}
	}
}
