package reserve

import (
	"fmt"
	"sync"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	reserves []Reserve
	sync.Mutex
}

func NewMemoryRepository(reserves *[]Reserve, filter Reserve, ordered bool) (mr MemoryRepository) {
	if reserves == nil {
		return MemoryRepository{
			reserves: []Reserve{},
		}
	}

	mr.reserves = []Reserve{}
	for _, res := range *reserves {
		if filter.Person.Id != uuid.Nil && res.Person.Id != filter.Person.Id {
			continue
		}

		if filter.StartTime != (time.Time{}) && filter.StartTime.After(res.EndTime) {
			continue
		}

		if filter.EndTime != (time.Time{}) && filter.EndTime.Before(res.StartTime) {
			continue
		}

		if filter.StartTime != (time.Time{}) && filter.EndTime != (time.Time{}) {
			fres := Reserve{
				StartTime: filter.StartTime,
				EndTime:   filter.EndTime,
			}
			if !fres.CheckConflicts(res) {
				continue
			}
		}
		if ordered {
			if !res.Ordered() {
				continue
			}
		}

		mr.reserves = append(mr.reserves, res)
	}
	return
}

func (mr *MemoryRepository) Get(id uuid.UUID) (Reserve, error) {
	for _, res := range mr.reserves {
		if res.Id == id {
			return res, nil
		}
	}
	return Reserve{}, ErrReserveNotFound
}

func (rep *MemoryRepository) GetByFilter(filter Reserve, ordered bool, sorted bool) (res []Reserve, err error) {
	newrep := NewMemoryRepository(&rep.reserves, filter, ordered)
	return newrep.reserves, nil
}

func (rep *MemoryRepository) Add(r Reserve) (res Reserve, err error) {
	if rep.reserves == nil {
		rep.Lock()
		rep.reserves = []Reserve{}
		rep.Unlock()
	}

	for _, rr := range rep.reserves {
		if rr.Id == r.Id {
			err = fmt.Errorf("reserve already exists: %w", ErrFailedToAddReserve)
			return
		}
	}
	rep.Lock()
	rep.reserves = append(rep.reserves, r)
	res = r
	rep.Unlock()
	return
}

func (mr *MemoryRepository) Update(memr Reserve) error {
	for idx, res := range mr.reserves {
		if res.Id == memr.Id {
			mr.Lock()
			mr.reserves[idx] = memr
			mr.Unlock()
			return nil
		}
	}
	return fmt.Errorf("reserve does not exist: %w", ErrUpdateReserve)
}

func (mr *MemoryRepository) AddPlayer(r Reserve, pl person.Person, count int) (Reserve, error) {
	for i, p := range r.Players {
		if p.Id == pl.Id {
			r.Players[i] = Player{Person: pl, Count: count}
			return r, nil
		}
	}
	r.Players = append(r.Players, Player{Person: pl, Count: count})
	return r, nil
}

func (mr *MemoryRepository) UpdatePlayer(r Reserve, pl person.Person, count int) (Reserve, error) {
	for i, p := range r.Players {
		if p.Id == pl.Id {
			r.Players[i] = Player{Person: pl, Count: count}
			return r, nil
		}
	}
	return r, ErrReservePlayerNotFound
}
