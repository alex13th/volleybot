package handlers

import (
	"fmt"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/telegram"

	"github.com/google/uuid"
)

type Updater interface {
	UpdateCQ(cq *telegram.CallbackQuery, state string) (telegram.MessageResponse, error)
	UpdateMsg(msg *telegram.Message) (telegram.MessageResponse, error)
}

type ReserveUpdater struct {
	CommonHandler
	PlayerHandler *PlayerHandler
	Reserve       reserve.Reserve
	Reserves      reserve.ReserveRepository
}

func NewReserveUpdater(data string, kh telegram.KeyboardHelper) (ru ReserveUpdater, err error) {
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
		ru.Reserve, err = ru.Reserves.Get(id)
		if err != nil {
			err = telegram.HelperError{
				Msg:       fmt.Sprintf("Getting reserve error: %s", err.Error()),
				AnswerMsg: "Getting reserve error"}
		}
	}
	return
}

func (ru *ReserveUpdater) UpdateCQ(cq *telegram.CallbackQuery, state string) (resp telegram.MessageResponse, err error) {
	st := telegram.State{
		ChatId:    cq.Message.Chat.Id,
		Data:      ru.Reserve.Id.String(),
		State:     state,
		MessageId: cq.Message.MessageId,
	}

	if err = ru.UpdateReserve(&ru.Reserve); err != nil {
		return ru.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	rm := NewReserveMessager(ru.Reserve, nil, ru.Resources)
	ru.StateRepository.Set(st)
	ru.UpdateReserveMessages(ru.Reserve, true)
	ru.PlayerHandler.NotifyPlayers(ru.Reserve, &rm, cq.From.Id, "notify")
	resp = cq.Answer(ru.Bot, "Ok", nil)
	return resp, nil
}

func (ru *ReserveUpdater) UpdateMsg(msg *telegram.Message) (resp telegram.MessageResponse, err error) {
	st := telegram.State{Data: ru.Reserve.Id.String(), State: "ordershow"}
	if err = ru.UpdateReserve(&ru.Reserve); err != nil {
		return ru.SendMessageError(msg, err.(telegram.HelperError), nil)
	}
	ru.StateRepository.Set(st)
	rm := NewReserveMessager(ru.Reserve, nil, ru.Resources)
	ru.UpdateReserveMessages(ru.Reserve, true)
	ru.PlayerHandler.NotifyPlayers(ru.Reserve, &rm, msg.From.Id, "notify")
	return
}

func (rh *ReserveUpdater) UpdateReserve(res *reserve.Reserve) (err error) {

	if err = rh.Reserves.Update(*res); err != nil {
		err = telegram.HelperError{
			Msg:       fmt.Sprintf("order update can't update reserve %s error: %s", res.Id, err.Error()),
			AnswerMsg: "Can't update reserve"}
		return
	}
	*res, err = rh.Reserves.Get(res.Id)
	return
}

func (rh *ReserveUpdater) UpdateReserveMessages(res reserve.Reserve, renew bool) {
	slist, _ := rh.StateRepository.GetByData(res.Id.String())
	notified := map[int]bool{}
	for _, st := range slist {
		if notified[st.ChatId] {
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
