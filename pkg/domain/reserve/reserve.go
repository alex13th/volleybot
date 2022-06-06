package reserve

import (
	"errors"
	"time"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

var (
	ErrReserveInvalidPeriod = errors.New("the reserve was not found in the repository")
	ErrFailedToAddReserve   = errors.New("failed to add the reserve to the repository")
	ErrUpdateReserve        = errors.New("failed to update the reserve in the repository")
	ErrReserveNotFound      = errors.New("a reserve has to have an valid person")
)

func NewPreReserve(p person.Person) Reserve {
	return Reserve{
		Id:         uuid.New(),
		Person:     p,
		CourtCount: 1,
		MaxPlayers: 6,
		Players:    make(map[uuid.UUID]Player)}
}

func NewReserve(p person.Person, start time.Time, end time.Time) (res Reserve, err error) {

	if end.Before(start) {
		return Reserve{}, ErrReserveInvalidPeriod
	}
	res = NewPreReserve(p)
	res.StartTime = start
	res.EndTime = end
	return
}

type Player struct {
	Person person.Person
	Count  int
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
	lnames[0] = "Не важен"
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

type Reserve struct {
	Id          uuid.UUID            `json:"id"`
	Person      person.Person        `json:"person"`
	Location    location.Location    `json:"location"`
	StartTime   time.Time            `json:"start_time"`
	EndTime     time.Time            `json:"end_time"`
	MinLevel    int                  `json:"min_level"`
	Price       int                  `json:"price"`
	CourtCount  int                  `json:"court_count"`
	MaxPlayers  int                  `json:"max_players"`
	Approved    bool                 `json:"approved"`
	Canceled    bool                 `json:"canceled"`
	Players     map[uuid.UUID]Player `json:"players"`
	Description string               `json:"description"`
}

func (res *Reserve) GetPerson() person.Person {
	return res.Person
}

func (res *Reserve) GetStartTime() time.Time {
	return res.StartTime
}

func (res *Reserve) GetEndTime() time.Time {
	return res.EndTime
}

func (res *Reserve) GetDuration() time.Duration {
	result := res.EndTime.Sub(res.StartTime)
	return result
}

func (res *Reserve) HasPlayerByTelegramId(id int) bool {
	for _, pl := range res.Players {
		if pl.Person.TelegramId == id {
			return true
		}
	}
	return false
}

func (res *Reserve) CheckConflicts(other Reserve) bool {

	OtherStartTime := other.GetStartTime()
	if res.StartTime == OtherStartTime {
		return true
	}

	if res.StartTime.Before(OtherStartTime) && OtherStartTime.Before(res.GetEndTime()) {
		return true
	}

	if res.StartTime.After(OtherStartTime) && res.StartTime.Before(other.GetEndTime()) {
		return true
	}

	return false
}

func (res Reserve) Orderd() (ordered bool) {
	ordered = (!res.StartTime.IsZero() && res.GetDuration() > 0 &&
		res.CourtCount > 0 && res.MaxPlayers > 0 && !res.Canceled)
	return
}
