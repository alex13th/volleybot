package handlers

import (
	"time"
	"volleybot/pkg/telegram"
)

type KeyboardCallbackHandler struct {
	KeyboardHelper telegram.KeyboardHelper
}

func (h *KeyboardCallbackHandler) Callback(cq *telegram.CallbackQuery) (result telegram.MessageResponse, err error) {
	kh := h.KeyboardHelper
	if err = kh.Parse(cq.Data); err != nil {
		return
	}

	if kh.Action == "set" {
		if err != nil {
			return h.SendCallbackError(cq, err.(telegram.HelperError), nil)
		}
		dur := res.GetDuration()
		res.StartTime = dh.Date.Add(time.Duration(res.StartTime.Hour()*int(time.Hour) +
			res.StartTime.Minute()*int(time.Minute)))
		res.EndTime = res.StartTime.Add(dur)

		return h.ReserveHandler.UpdateReserveCQ(res, cq, "ordershow", false)
	} else {
		rm := NewReserveMessager(res, &dh, h.Resources)
		mr := rm.GetEditMR(cq.Message.Chat.Id)
		cq.Message.EditText(h.Bot, "", &mr)
		return cq.Answer(h.Bot, h.Resources.OkAnswer, nil), nil
	}
}
