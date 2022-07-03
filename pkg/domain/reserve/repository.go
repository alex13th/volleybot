package reserve

import (
	"volleybot/pkg/domain/person"

	uuid "github.com/google/uuid"
)

type ReserveRepository interface {
	Add(Reserve) (Reserve, error)
	AddPlayer(Reserve, person.Player) (Reserve, error)
	Get(uuid.UUID) (Reserve, error)
	GetByFilter(res Reserve, oredered bool, sorted bool) ([]Reserve, error)
	UpdatePlayer(Reserve, person.Player) (Reserve, error)
	Update(Reserve) error
}
