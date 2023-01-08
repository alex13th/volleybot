package reserve

import (
	"encoding/base64"
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
		Id:     uuid.New(),
		Person: p,
	}
}

func NewReserve(p person.Person, start time.Time, end time.Time) (res Reserve) {
	res = NewPreReserve(p)
	res.StartTime = start
	res.EndTime = end
	return
}

type Reserve struct {
	Id          uuid.UUID         `json:"id"`
	Person      person.Person     `json:"person"`
	Location    location.Location `json:"location"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Price       int               `json:"price"`
	Approved    bool              `json:"approved"`
	Canceled    bool              `json:"canceled"`
	Description string            `json:"description"`
}

func (res Reserve) Base64Id() string {
	bid := [16]byte(res.Id)
	return base64.StdEncoding.EncodeToString(bid[:])
}

func (res Reserve) IdFromBase64(b64 string) (id uuid.UUID, err error) {
	var bid []byte
	if bid, err = base64.StdEncoding.DecodeString(b64); err != nil {
		return
	}
	id, err = uuid.FromBytes(bid)
	return
}

func (res *Reserve) GetPerson() person.Person {
	return res.Person
}

func (res *Reserve) GetStartTime() time.Time {
	return res.StartTime
}

func (res *Reserve) SetDurationHours(h int) {
	res.EndTime = res.StartTime.Add(time.Duration(time.Hour * time.Duration(h)))
}

func (res *Reserve) SetStartDate(dt time.Time) {
	dur := res.GetDuration()
	res.StartTime = dt.Add(time.Duration(res.StartTime.Hour()*int(time.Hour) +
		res.StartTime.Minute()*int(time.Minute)))
	res.EndTime = res.StartTime.Add(dur)
}

func (res *Reserve) SetStartTime(tm time.Time) {
	dur := res.GetDuration()
	res.StartTime = time.Date(res.StartTime.Year(), res.StartTime.Month(), res.StartTime.Day(),
		tm.Hour(), tm.Minute(), 0, 0, time.Local)
	res.EndTime = res.StartTime.Add(dur)
}

func (res *Reserve) GetEndTime() time.Time {
	return res.EndTime
}

func (res *Reserve) GetDuration() time.Duration {
	result := res.EndTime.Sub(res.StartTime)
	return result
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
		!res.Canceled)
	return
}
