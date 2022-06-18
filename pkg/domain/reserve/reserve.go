package reserve

import (
	"errors"
	"time"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

var (
	ErrReserveInvalidPeriod  = errors.New("the reserve was not found in the repository")
	ErrFailedToAddReserve    = errors.New("failed to add the reserve to the repository")
	ErrUpdateReserve         = errors.New("failed to update the reserve in the repository")
	ErrReserveNotFound       = errors.New("a reserve has to have an valid person")
	ErrReservePlayerNotFound = errors.New("a reserve has to have an valid player")
)

func NewPreReserve(p person.Person) Reserve {
	return Reserve{
		Id:         uuid.New(),
		Person:     p,
		CourtCount: 1,
		MaxPlayers: 6,
		Players:    []person.Player{}}
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

type Activity int

const (
	Game       Activity = 0
	Training   Activity = 10
	Tournament Activity = 20
)

func (a Activity) String() string {
	lnames := make(map[int]string)
	lnames[0] = "🏐 Игры"
	lnames[10] = "‼️ Тренировка"
	lnames[20] = "🌟 Турнир"
	return lnames[int(a)]
}

type Reserve struct {
	Id          uuid.UUID         `json:"id"`
	Activity    int               `json:"activity"`
	Person      person.Person     `json:"person"`
	Location    location.Location `json:"location"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	MinLevel    int               `json:"min_level"`
	Price       int               `json:"price"`
	CourtCount  int               `json:"court_count"`
	MaxPlayers  int               `json:"max_players"`
	Approved    bool              `json:"approved"`
	Canceled    bool              `json:"canceled"`
	Players     []person.Player   `json:"players"`
	Description string            `json:"description"`
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

func (res *Reserve) Copy() (result Reserve) {
	result = *res
	result.Id = uuid.New()
	return
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

func (res Reserve) Ordered() (ordered bool) {
	ordered = (!res.StartTime.IsZero() && res.GetDuration() > 0 &&
		res.CourtCount > 0 && res.MaxPlayers > 0 && !res.Canceled)
	return
}

func (res *Reserve) PlayerCount(pid uuid.UUID) (count int) {
	for i, pl := range res.Players {
		if res.Players[i].Id != pid {
			count += pl.Count
		}
	}
	return
}

func (res *Reserve) GetPlayer(pid uuid.UUID) (pl person.Player) {
	for _, pl := range res.Players {
		if pl.Id != pid {
			return pl
		}
	}
	return
}

func (res *Reserve) JoinPlayer(pl person.Player) {
	for i, p := range res.Players {
		if p.Id == pl.Id {
			res.Players[i] = pl
		}
	}
	res.Players = append(res.Players, pl)
}
