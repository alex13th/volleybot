package person

import (
	uuid "github.com/google/uuid"
)

type PersonRepository interface {
	Get(uuid.UUID) (Person, error)
	Add(Person) error
	Update(Person) error
}
