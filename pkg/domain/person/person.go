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

	Params        = []string{"notify", "notify_cancel"}
	ParamDefaults = map[string]string{
		"notify":        "off",
		"notify_cancel": "off",
	}
	ParamValText = map[string]string{
		"undef": "не определен",
		"on":    "вкл.",
		"off":   "выкл.",
	}

	ParamNames = map[string]string{
		"notify":        "При изменении",
		"notify_cancel": "При отмене",
	}
)

func NewPerson(firstname string) Person {
	return Person{
		Firstname:     firstname,
		Id:            uuid.New(),
		LocationRoles: make(map[uuid.UUID][]string),
		Settings:      make(map[string]string),
	}
}

type Person struct {
	Id            uuid.UUID              `json:"id"`
	TelegramId    int                    `json:"telegram_id"`
	Firstname     string                 `json:"firstname"`
	Lastname      string                 `json:"lastname"`
	Fullname      string                 `json:"fullname"`
	Sex           Sex                    `json:"sex"`
	LocationRoles map[uuid.UUID][]string `json:"roles"`
	Settings      map[string]string      `json:"settings"`
}

func (user *Person) String() string {
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

type Sex int

func (s Sex) String() string {
	lnames := make(map[int]string)
	lnames[0] = ""
	lnames[1] = "мальчик"
	lnames[2] = "девочка"
	return lnames[int(s)]
}

func (s Sex) Emoji() string {
	lnames := make(map[int]string)
	lnames[0] = "👤"
	lnames[1] = "👦🏻"
	lnames[2] = "👩🏻"
	return lnames[int(s)]
}
