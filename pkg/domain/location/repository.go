package location

import (
	uuid "github.com/google/uuid"
)

type LocationRepository interface {
	Get(uuid.UUID) (Location, error)
	Add(Location) (Location, error)
	Update(Location) error
}
