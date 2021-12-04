package reserve

import (
	"time"
	"volleybot/pkg/domain/person"

	uuid "github.com/google/uuid"
)

type ReserveRepository interface {
	Get(uuid.UUID) (Reserve, error)
	GetByFilter(ReserveFilter) (map[uuid.UUID]Reserve, error)
	Add(Reserve) error
	Update(Reserve) error
}

type ReserveFilter struct {
	Person    person.Person
	StartTime time.Time
	EndTime   time.Time
}
