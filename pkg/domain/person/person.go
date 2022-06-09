package person

import (
	"errors"
	"strings"
	"volleybot/pkg/domain/location"

	uuid "github.com/google/uuid"
)

type ErrorPersonNotFound struct {
	msg string
}

func (e ErrorPersonNotFound) Error() string {
	return e.msg
}

var (
	ErrPersonNotFound    = ErrorPersonNotFound{msg: "the person was not found in the repository"}
	ErrFailedToAddPerson = errors.New("failed to add the person to the repository")
	ErrUpdatePerson      = errors.New("failed to update the person in the repository")
	ErrInvalidPerson     = errors.New("a person has to have an valid name")
)

func NewPerson(firstname string) (person Person, err error) {
	if firstname == "" {
		return Person{}, ErrInvalidPerson
	}

	person = Person{
		Firstname: firstname,
		Id:        uuid.New(),
	}
	return
}

type Person struct {
	Id            uuid.UUID              `json:"id"`
	TelegramId    int                    `json:"telegram_id"`
	Firstname     string                 `json:"firstname"`
	Lastname      string                 `json:"lastname"`
	Fullname      string                 `json:"fullname"`
	LocationRoles map[uuid.UUID][]string `json:"roles"`
	Settings      map[string]string      `json:"settings"`
}

func (user *Person) GetDisplayname() string {
	fullname := strings.Trim(user.Fullname, " ")
	if fullname != "" {
		return fullname
	}
	firstname := strings.Trim(user.Firstname, " ")
	lastname := strings.Trim(user.Lastname, " ")

	if lastname != "" {
		if firstname != "" {
			return firstname + " " + lastname
		} else {
			return lastname
		}
	}

	return firstname
}

func (user *Person) CheckLocationRole(l location.Location, role string) bool {
	for _, r := range user.LocationRoles[l.Id] {
		if r == role {
			return true
		}
	}
	return false
}
