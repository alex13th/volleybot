package reserve

import (
	uuid "github.com/google/uuid"
)

type ReserveRepository interface {
	Add(Reserve) (Reserve, error)
	Get(uuid.UUID) (Reserve, error)
	GetByFilter(res Reserve, oredered bool, sorted bool) ([]Reserve, error)
	Update(Reserve) error
}
