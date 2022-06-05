package location

import (
	"errors"

	uuid "github.com/google/uuid"
)

var (
	ErrLocationNotFound  = errors.New("the location was not found in the repository")
	ErrFailedToAddPerson = errors.New("failed to add the location to the repository")
	ErrUpdatePerson      = errors.New("failed to update the location in the repository")
	ErrInvalidPerson     = errors.New("a location has to have an valid name")
)

func NewLocation(name string) (location Location, err error) {
	if name == "" {
		return Location{}, ErrInvalidPerson
	}

	location = Location{
		Name: name,
		Id:   uuid.New(),
	}
	return
}

type Location struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ChatId      int       `json:"chat_id"`
}
