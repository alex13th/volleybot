package reserve

import (
	uuid "github.com/google/uuid"
)

type ReserveRepository interface {
	Get(uuid.UUID) (Reserve, error)
	Add(Reserve) error
	Update(Reserve) error
}
