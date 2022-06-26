package telegram

import (
	"errors"
	"sync"
)

var (
	ErrStateNotFound    = errors.New("the state was not found in the repository")
	ErrFailedToAddState = errors.New("failed to add the state to the repository")
	ErrUpdateState      = errors.New("failed to update the state in the repository")
)

type State struct {
	ChatId    int    `json:"chat_id"`
	MessageId int    `json:"message_id"`
	State     string `json:"state"`
	Data      string `json:"data"`
}

type StateRepository interface {
	Get(ChatId int) ([]State, error)
	GetByData(Data string) ([]State, error)
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
