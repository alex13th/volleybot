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
		"notify_cancel": "on",
	}
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
	Sex           Sex                    `json:"sex"`
	Level         PlayerLevel            `json:"level"`
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

type Player struct {
	Person
	Count int
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

type PlayerLevel int

const (
	Nothing      PlayerLevel = 0
	Novice       PlayerLevel = 10
	Begginer     PlayerLevel = 20
	BeginnerPlus PlayerLevel = 30
	MiddleMinus  PlayerLevel = 40
	Middle       PlayerLevel = 50
	MiddlePlus   PlayerLevel = 60
	Advanced     PlayerLevel = 70
	Proffesional PlayerLevel = 80
)

func (pl PlayerLevel) String() string {
	lnames := make(map[int]string)
	lnames[0] = "Не определен"
	lnames[10] = "Новичок"
	lnames[20] = "Начальный"
	lnames[30] = "Начальный+"
	lnames[40] = "Средний-"
	lnames[50] = "Средний"
	lnames[60] = "Средний+"
	lnames[70] = "Уверенный"
	lnames[80] = "Профи"
	return lnames[int(pl)]
}

func (pl PlayerLevel) Emoji() string {
	lnames := make(map[int]string)
	lnames[0] = ""
	lnames[10] = "🙌"
	lnames[20] = "👏"
	lnames[30] = "🤝"
	lnames[40] = "👌"
	lnames[50] = "👍"
	lnames[60] = "💪"
	lnames[70] = "⭐️"
	lnames[80] = "👑"
	return lnames[int(pl)]
}
