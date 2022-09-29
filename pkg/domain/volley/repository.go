package volley

import (
	"volleybot/pkg/domain/person"

	uuid "github.com/google/uuid"
)

type Repository interface {
	Add(Volley) (Volley, error)
	AddMember(Volley, Member) (Volley, error)
	AddPlayer(Player) (Player, error)
	Get(uuid.UUID) (Volley, error)
	GetByFilter(res Volley, oredered bool, sorted bool) ([]Volley, error)
	GetPlayer(person.Person) (Player, error)
	UpdateMember(Volley, Member) (Volley, error)
	UpdatePlayer(Player) error
	Update(Volley) error
}
