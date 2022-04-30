package person

import (
	uuid "github.com/google/uuid"
)

type PersonRepository interface {
	Get(uuid.UUID) (Person, error)
	GetByTelegramId(int) (Person, error)
	Add(Person) (Person, error)
	Update(Person) error
}
