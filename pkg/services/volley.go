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
	bp, err := p.GetBaseStateProvider(cq.From.Id, st, *cq.Message)
	if err != nil {
		p.LogErrors([]error{err})
		return
	}
	p.LogErrors(p.Proceed(bp))
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
		st.Prefix = "res"
		bp, err := p.GetBaseStateProvider(msg.From.Id, st, *msg)
		if err != nil {
			p.LogErrors([]error{err})
			return err
		}
		p.LogErrors(p.Proceed(bp))
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
	bp, err := p.GetBaseStateProvider(msg.From.Id, st, stmsg)
	if err != nil {
		p.LogErrors([]error{err})
		return
	}
	p.LogErrors(p.Proceed(bp))
	return
}

func (p *VolleyBotService) Proceed(bp bvbot.BaseStateProvider) (errs []error) {
	var (
		err     error
		reqlist []telegram.StateRequest
		sp      telegram.StateProvider
		state   telegram.State
	)
	if sp, err = p.GetStateProvider(bp); sp == nil {
		return append(errs, err)
	}
	if state, err = sp.Proceed(); sp == nil {
		return append(errs, err)
	}
	// Adding incoming state requests
	reqlist = append(reqlist, sp.GetRequests()...)
	errs = append(errs, p.SendRequests(reqlist)...)
	reqlist = []telegram.StateRequest{}

	// Adding result state requests
	bp, err = bvbot.NewBaseStateProvider(state, bp.Message, bp.Person, bp.Location, bp.Repository,
		p.Resources.Volley.MaxPlayer.GroupChatWarning)
	if err != nil {
		return append(errs, err)
	}
	if state.Updated {
		p.UpdateMessages(bp)
	}
	if sp, err = p.GetStateProvider(bp); sp == nil {
		errs = append(errs, err)
	} else {
		reqlist = append(reqlist, sp.GetRequests()...)
	}

	return append(errs, p.SendRequests(reqlist)...)
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

func (s *VolleyBotService) GetBaseStateProvider(tid int, state telegram.State, msg telegram.Message) (bp bvbot.BaseStateProvider, err error) {
	p, err := s.PersonRepository.GetByTelegramId(tid)
	if err != nil {
		p, _ = person.NewPerson(msg.From.FirstName)
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

	bp, err = bvbot.NewBaseStateProvider(state, msg, p, loc, s.VolleyRepository,
		s.Resources.Volley.MaxPlayer.GroupChatWarning)
	if err != nil {
		return
	}
	return
}

func (p *VolleyBotService) GetStateProvider(bp bvbot.BaseStateProvider) (sp telegram.StateProvider, err error) {
	bp.BackState = bp.State
	switch bp.State.State {
	case "main":
		bp.BackState = telegram.State{}
		sp = bvbot.MainStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Main}
	case "listd":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.ListdStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.List}
	case "show":
		bp.BackState.State = "main"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.ShowStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Show}
	case "actions":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.ActionsStateProvider{BaseStateProvider: bp,
			Resources: p.Resources.Volley.Actions, ShowResources: p.Resources.Volley.Show}
	case "date":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.DateStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Show.DateTime}
	case "desc":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = &bvbot.DescStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Description}
	case "time":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.TimeStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Show.DateTime}
	case "sets":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.SetsStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Sets}
	case "joinm":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.JoinPlayersStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Join}
	case "jtime":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.JoinTimeStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Join}
	case "settings":
		bp.BackState.State = "show"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.SettingsStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Settings}
	case "courts":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.CourtsStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Courts}
	case "max":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.MaxPlayersStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.MaxPlayer}
	case "price":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.PriceStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Price}
	case "level":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.LevelStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Level}
	case "activity":
		bp.BackState.State = "settings"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.ActivityStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.Activity}
	case "cancel":
		bp.BackState.State = "actions"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.CancelStateProvider{BaseStateProvider: bp,
			Resources: p.Resources.Volley.Cancel, ShowResources: p.Resources.Volley.Show}
	case "rmpl":
		bp.BackState.State = "actions"
		bp.BackState.Action = bp.BackState.State
		sp = bvbot.RemovePlayerStateProvider{BaseStateProvider: bp, Resources: p.Resources.Volley.RemovePlayer}
	}
	return
}

func (p *VolleyBotService) UpdateMessages(bp bvbot.BaseStateProvider) {
	slist, _ := p.StateRepository.GetByData(bp.State.Data)
	sort.Slice(slist, func(i, j int) bool {
		return slist[i].MessageId > slist[j].MessageId
	})
	cid := bp.Message.Chat.Id
	mid := bp.Message.MessageId
	bp.Message = telegram.Message{Chat: &telegram.Chat{Id: bp.Message.Chat.Id}}
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

		bp.Message.Chat.Id = st.ChatId
		bp.Message.MessageId = st.MessageId
		bp.State = st
		if sp, _ = p.GetStateProvider(bp); sp != nil {
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
