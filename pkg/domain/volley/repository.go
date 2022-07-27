package volley

import (
	"volleybot/pkg/domain/person"

	uuid "github.com/google/uuid"
)

type Repository interface {
	Add(Volley) (Volley, error)
	AddPlayer(Volley, person.Player) (Volley, error)
	Get(uuid.UUID) (Volley, error)
	GetByFilter(res Volley, oredered bool, sorted bool) ([]Volley, error)
	UpdatePlayer(Volley, person.Player) (Volley, error)
	Update(Volley) error
}
