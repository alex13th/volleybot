package location

import (
	uuid "github.com/google/uuid"
)

type LocationRepository interface {
	Get(uuid.UUID) (Location, error)
	GetByName(string) (Location, error)
	Add(Location) (Location, error)
	Update(Location) error
}

type LocationConfigRepository interface {
	Add(loc Location, service string, config interface{}) error
	Get(loc Location, service string, config interface{}) error
	Update(loc Location, service string, config interface{}) error
}
