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
	Get(ChatId int) (State, error)
	Set(State) error
	Clear(ChatId int)
}

func NewMemoryStateRepository() StateRepository {
	return &MemoryStateRepository{states: make(map[int]State)}
}

type MemoryStateRepository struct {
	states map[int]State
	sync.Mutex
}

func (rep *MemoryStateRepository) Get(ChatId int) (st State, err error) {
	if state, ok := rep.states[ChatId]; ok {
		return state, nil
	}
	return State{}, ErrStateNotFound
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

func (rep *MemoryStateRepository) Clear(ChatId int) {
	delete(rep.states, ChatId)
}
