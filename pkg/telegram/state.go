package telegram

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrStateNotFound    = errors.New("the state was not found in the repository")
	ErrFailedToAddState = errors.New("failed to add the state to the repository")
	ErrUpdateState      = errors.New("failed to update the state in the repository")
)

type StateBuilder interface {
	GetStateProvider(State) (StateProvider, error)
}

type StateProvider interface {
	GetRequests() []StateRequest
	Proceed() (State, error)
}

func NewState() State {
	return State{Separator: "_"}
}

type State struct {
	Action    string `json:"action"`
	ChatId    int    `json:"chat_id"`
	Data      string `json:"data"`
	MessageId int    `json:"message_id"`
	Prefix    string `json:"prefix"`
	Separator string `json:"separator"`
	State     string `json:"state"`
	Updated   bool   `json:"updated"`
	Value     string `json:"value"`
}

func (st State) String() string {
	if st.Action == "" {
		st.Action = st.State
	}
	slist := []string{st.Prefix, st.State, st.Action}
	if st.Data != "" {
		slist = append(slist, st.Data)
	}
	if st.Value != "" {
		slist = append(slist, st.Value)
	}
	return strings.Join(slist, st.Separator)
}

func (st State) Parse(data string) (state State, err error) {
	state = st
	splitedData := strings.Split(data, st.Separator)
	if len(splitedData) < 3 {
		err = HelperError{
			Msg:       "incorrect Date button data format",
			AnswerMsg: "Can't parse date"}
		return
	}
	state.Prefix = splitedData[0]
	state.State = splitedData[1]
	state.Action = splitedData[2]
	if len(splitedData) > 3 {
		state.Data = splitedData[3]
	}
	if len(splitedData) > 4 {
		state.Value = strings.Join(splitedData[4:], st.Separator)
	}
	return
}

type StateRequest struct {
	State
	Request
	Clear bool
}

type StateRepository interface {
	Get(ChatId int) ([]State, error)
	GetByData(Data string) ([]State, error)
	GetByMessage(msg Message) (State, error)
	Set(State) error
	Clear(State) error
}

func NewMemoryStateRepository() StateRepository {
	return &MemoryStateRepository{states: make(map[int]State)}
}

type MemoryStateRepository struct {
	states map[int]State
	sync.Mutex
}

func (rep *MemoryStateRepository) Get(ChatId int) (st []State, err error) {
	if state, ok := rep.states[ChatId]; ok {
		return []State{state}, nil
	}
	return []State{}, ErrStateNotFound
}

func (rep *MemoryStateRepository) GetByMessage(msg Message) (st State, err error) {
	if state, ok := rep.states[msg.Chat.Id]; ok {
		return state, nil
	}
	return State{}, ErrStateNotFound
}

func (rep *MemoryStateRepository) GetByData(Data string) (slist []State, err error) {
	rep.Mutex.Lock()
	for _, s := range rep.states {
		if s.Data == Data {
			slist = append(slist, s)
		}
	}
	rep.Mutex.Unlock()
	return
}

func (rep *MemoryStateRepository) Set(s State) (err error) {
	if rep.states == nil {
		rep.Lock()
		rep.states = make(map[int]State)
		rep.Unlock()
	}
	rep.Lock()
	rep.states[s.ChatId] = s
	rep.Unlock()
	return
}

func (rep *MemoryStateRepository) Clear(st State) error {
	delete(rep.states, st.ChatId)
	return nil
}
