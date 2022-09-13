package services

import (
	"log"
	"sort"
	"volleybot/pkg/bvbot"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/volley"
	"volleybot/pkg/res"
	"volleybot/pkg/telegram"
)

type VolleyBotService struct {
	Bot                *telegram.Bot
	Prefix             string
	Resources          *res.VolleyResources
	LocationRepository location.LocationRepository
	PersonRepository   person.PersonRepository
	VolleyRepository   volley.Repository
	StateRepository    telegram.StateRepository
}

func (s VolleyBotService) LogErrors(errs []error) {
	for _, err := range errs {
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (s *VolleyBotService) GetCommands(tuser int) (cmds []telegram.BotCommand) {
	cmds = append(cmds, s.Resources.Command)
	return cmds
}

func (s *VolleyBotService) GetLocation() (l location.Location, err error) {
	l, err = s.LocationRepository.GetByName(s.Resources.Location.Name)
	if err != nil {
		log.Println(err.Error())
		l, _ = location.NewLocation(s.Resources.Location.Name)
		l, err = s.LocationRepository.Add(l)
	}
	return
}

func (p *VolleyBotService) ProceedCallback(cq *telegram.CallbackQuery) (err error) {
	st, err := telegram.NewState().Parse(cq.Data)
	st.MessageId = cq.Message.MessageId
	st.ChatId = cq.Message.Chat.Id
	if err != nil {
		log.Println(err.Error())
		return
	}
	p.LogErrors(p.Proceed(cq.From.Id, st, *cq.Message))
	_, err = cq.Answer(p.Bot, "Ok", nil)
	return
}

func (p *VolleyBotService) ProceedMessage(msg *telegram.Message) (err error) {
	cmd := msg.GetCommand()
	switch cmd {
	case "volley":
		st := telegram.NewState()
		st.Action = "start"
		st.State = "main"
		st.ChatId = msg.Chat.Id
		st.Prefix = "res"
		p.LogErrors(p.Proceed(msg.From.Id, st, *msg))
		return err
	}

	slist, err := p.StateRepository.Get(msg.Chat.Id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(slist) == 0 {
		return
	}
	st := slist[0]
	stmsg := *msg
	stmsg.MessageId = st.MessageId
	p.LogErrors(p.Proceed(msg.From.Id, st, stmsg))
	return
}

func (p *VolleyBotService) Proceed(tid int, st telegram.State, msg telegram.Message) (errs []error) {
	var (
		err      error
		reqlist  []telegram.StateRequest
		sp       telegram.StateProvider
		newstate telegram.State
	)
	bld, err := p.GetStateBuilder(tid, st, msg)
	if err != nil {
		return append(errs, err)
	}

	if sp, err = bld.GetStateProvider(st); sp == nil {
		return append(errs, err)
	}
	if newstate, err = sp.Proceed(); sp == nil {
		return append(errs, err)
	}
	// Adding incoming state requests
	reqlist = append(reqlist, sp.GetRequests()...)
	errs = append(errs, p.SendRequests(reqlist)...)
	reqlist = []telegram.StateRequest{}

	// Adding result state requests
	bld, err = p.GetStateBuilder(tid, newstate, msg)
	if err != nil {
		return append(errs, err)
	}
	if sp, err = bld.GetStateProvider(newstate); sp == nil {
		errs = append(errs, err)
	} else {
		reqlist = append(reqlist, sp.GetRequests()...)
	}

	errs = append(errs, p.SendRequests(reqlist)...)

	if newstate.Updated {
		p.UpdateMessages(newstate, bld)
	}

	return
}

func (s *VolleyBotService) SendRequests(reqlist []telegram.StateRequest) (errs []error) {
	var err error
	for _, req := range reqlist {
		if req.Clear {
			if err = s.StateRepository.Clear(req.State); err != nil {
				errs = append(errs, err)
			}
			continue
		}
		var resp telegram.MessageResponse
		if resp, err = s.Bot.SendMessage(req.Request); err != nil {
			errs = append(errs, err)
			continue
		}
		if resp.Result.Chat != nil {
			if req.State.MessageId >= 0 {
				req.State.MessageId = resp.Result.MessageId
			}
			req.State.ChatId = resp.Result.Chat.Id
			if err = s.StateRepository.Set(req.State); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return
}

func (s *VolleyBotService) GetStateBuilder(tid int, state telegram.State, msg telegram.Message) (bld telegram.StateBuilder, err error) {
	p, err := s.PersonRepository.GetByTelegramId(tid)
	if err != nil {
		p = person.NewPerson(msg.From.FirstName)
		p.TelegramId = msg.From.Id
		p.Lastname = msg.From.LastName
		if p, err = s.PersonRepository.Add(p); err != nil {
			return
		}
	}
	loc, err := s.GetLocation()
	if err != nil {
		return
	}

	return bvbot.NewBvStateBuilder(loc, msg, p, s.VolleyRepository, s.Resources.Volley, state)
}

func (p *VolleyBotService) UpdateMessages(sta telegram.State, bld telegram.StateBuilder) {
	slist, _ := p.StateRepository.GetByData(sta.Data)
	sort.Slice(slist, func(i, j int) bool {
		return slist[i].MessageId > slist[j].MessageId
	})
	cid := sta.ChatId
	mid := sta.MessageId
	notified := map[int]bool{}
	for _, st := range slist {
		var (
			reqlist []telegram.StateRequest
			sp      telegram.StateProvider
		)
		if st.ChatId == cid && (st.MessageId == mid || st.ChatId < 0) {
			continue
		}

		if st.ChatId < 0 {
			p.StateRepository.Clear(st)
		}
		if notified[st.ChatId] {
			continue
		}
		notified[st.ChatId] = true

		if sp, _ = bld.GetStateProvider(st); sp != nil {
			reqlist = append(reqlist, sp.GetRequests()...)
		}
		for _, req := range reqlist {
			resp, err := p.Bot.SendMessage(req.Request)
			if err == nil {
				st.MessageId = resp.Result.MessageId
			}
		}
		p.StateRepository.Set(st)
	}
}
