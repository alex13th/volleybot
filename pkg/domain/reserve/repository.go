package reserve

import (
	"volleybot/pkg/domain/person"

	uuid "github.com/google/uuid"
)

type ReserveRepository interface {
	Get(uuid.UUID) (Reserve, error)
	GetByFilter(Reserve, bool) (map[uuid.UUID]Reserve, error)
	Add(Reserve) (Reserve, error)
	AddPlayer(Reserve, person.Person, int) (Reserve, error)
	UpdatePlayer(Reserve, person.Person, int) (Reserve, error)
	Update(Reserve) error
}
