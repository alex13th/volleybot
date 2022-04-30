package person

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	persons map[uuid.UUID]Person
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		persons: make(map[uuid.UUID]Person),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (Person, error) {
	if person, ok := mr.persons[id]; ok {
		return person, nil
	}

	return Person{}, ErrPersonNotFound
}

func (mr *MemoryRepository) GetByTelegramId(id int) (Person, error) {
	for _, person := range mr.persons {
		if person.TelegramId == id {
			return person, nil
		}
	}

	return Person{}, ErrPersonNotFound
}

func (mr *MemoryRepository) Add(p Person) (per Person, err error) {
	if mr.persons == nil {
		mr.Lock()
		mr.persons = make(map[uuid.UUID]Person)
		mr.Unlock()
	}

	if _, ok := mr.persons[p.Id]; ok {
		err = fmt.Errorf("person already exists: %w", ErrFailedToAddPerson)
		return
	}
	mr.Lock()
	mr.persons[p.Id] = p
	per = p
	mr.Unlock()
	return
}

func (mr *MemoryRepository) Update(memp Person) error {
	if _, ok := mr.persons[memp.Id]; !ok {
		return fmt.Errorf("person does not exist: %w", ErrUpdatePerson)
	}
	mr.Lock()
	mr.persons[memp.Id] = memp
	mr.Unlock()
	return nil
}
