package memory

import (
	"fmt"
	"sync"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	persons map[uuid.UUID]person.Person
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		persons: make(map[uuid.UUID]person.Person),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (person.Person, error) {
	if person, ok := mr.persons[id]; ok {
		return person, nil
	}

	return person.Person{}, person.ErrPersonNotFound
}

func (mr *MemoryRepository) Add(memp person.Person) error {
	if mr.persons == nil {
		mr.Lock()
		mr.persons = make(map[uuid.UUID]person.Person)
		mr.Unlock()
	}

	if _, ok := mr.persons[memp.Id]; ok {
		return fmt.Errorf("person already exists: %w", person.ErrFailedToAddPerson)
	}
	mr.Lock()
	mr.persons[memp.Id] = memp
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(memp person.Person) error {
	if _, ok := mr.persons[memp.Id]; !ok {
		return fmt.Errorf("person does not exist: %w", person.ErrUpdatePerson)
	}
	mr.Lock()
	mr.persons[memp.Id] = memp
	mr.Unlock()
	return nil
}
