package handlers

import (
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"
	"volleybot/pkg/services"
	"volleybot/pkg/telegram"
)

type PlayerHandler struct {
	CommonHandler
	PersonService *services.PersonService
}

func (ph *PlayerHandler) NotifyPlayers(res reserve.Reserve, rm Messager, tid int, action string) {
	for _, pl := range res.Players {
		if pl.Person.TelegramId != tid {
			p, _ := ph.PersonService.GetByTelegramId(pl.Person.TelegramId)
			if param, ok := p.Settings[action]; ok && param == "on" {
				mr := rm.GetMR(pl.TelegramId)
				ph.Bot.SendMessage(&mr)
				return
			}
		}
	}
}

func (ph *PlayerHandler) JoinPlayer(cq *telegram.CallbackQuery, res *reserve.Reserve, count int) (result telegram.MessageResponse, err error) {
	p, err := ph.GetPerson(cq.From)
	if err != nil {
		return ph.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if err != nil {
		return ph.SendCallbackError(cq, err.(telegram.HelperError), nil)
	}

	if res.HasPlayerByTelegramId(p.TelegramId) && res.MaxPlayers-count+res.PlayerCount(p.Id) < 0 {
		err := telegram.HelperError{Msg: "order join max players count error", AnswerMsg: ph.Resources.MaxPlayer.CountError}
		return ph.SendCallbackError(cq, err, nil)
	}
	res.JoinPlayer(person.Player{Person: p, Count: count})
	return
}
