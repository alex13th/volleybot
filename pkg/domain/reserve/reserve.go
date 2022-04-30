package reserve

import (
	"errors"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

var (
	ErrReserveInvalidPeriod = errors.New("the reserve was not found in the repository")
	ErrFailedToAddReserve   = errors.New("failed to add the reserve to the repository")
	ErrUpdateReserve        = errors.New("failed to update the reserve in the repository")
	ErrReserveNotFound      = errors.New("a reserve has to have an valid person")
)

func NewPreReserve(p person.Person) (res Reserve) {
	return Reserve{Id: uuid.New(), Person: p, Players: make(map[uuid.UUID]Player)}
}

func NewReserve(p person.Person, start time.Time, end time.Time) (reserve Reserve, err error) {

	if end.Before(start) {
		return Reserve{}, ErrReserveInvalidPeriod
	}

	reserve = Reserve{
		Id:        uuid.New(),
		Person:    p,
		StartTime: start,
		EndTime:   end,
		Players:   make(map[uuid.UUID]Player)}
	return
}

type Player struct {
	Person person.Person
	Count  int
}

type Reserve struct {
	Id         uuid.UUID            `json:"id"`
	Person     person.Person        `json:"person"`
	StartTime  time.Time            `json:"start_time"`
	EndTime    time.Time            `json:"end_time"`
	CourtCount int                  `json:"court_count"`
	MaxPlayers int                  `json:"max_players"`
	Ordered    bool                 `json:"ordered"`
	Approved   bool                 `json:"approved"`
	Players    map[uuid.UUID]Player `json:"players"`
}

func (reserve *Reserve) GetPerson() person.Person {
	return reserve.Person
}

func (reserve *Reserve) GetStartTime() time.Time {
	return reserve.StartTime
}

func (reserve *Reserve) GetEndTime() time.Time {
	return reserve.EndTime
}

func (reserve *Reserve) GetDuration() time.Duration {
	result := reserve.EndTime.Sub(reserve.StartTime)
	return result
}

func (reserve *Reserve) CheckConflicts(other Reserve) bool {

	OtherStartTime := other.GetStartTime()
	if reserve.StartTime == OtherStartTime {
		return true
	}

	if reserve.StartTime.Before(OtherStartTime) && OtherStartTime.Before(reserve.GetEndTime()) {
		return true
	}

	if reserve.StartTime.After(OtherStartTime) && reserve.StartTime.Before(other.GetEndTime()) {
		return true
	}

	return false
}
