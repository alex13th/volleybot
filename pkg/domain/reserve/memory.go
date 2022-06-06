package reserve

import (
	"fmt"
	"sync"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	reserves map[uuid.UUID]Reserve
	sync.Mutex
}

func NewMemoryRepository(reserves *map[uuid.UUID]Reserve, filter Reserve, ordered bool) (mr MemoryRepository) {
	if reserves == nil {
		return MemoryRepository{
			reserves: make(map[uuid.UUID]Reserve),
		}
	}

	mr.reserves = make(map[uuid.UUID]Reserve)
	for id, res := range *reserves {
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
			if !res.Orderd() {
				continue
			}
		}

		mr.reserves[id] = res
	}
	return
}

func (mr *MemoryRepository) Get(id uuid.UUID) (Reserve, error) {
	if reserve, ok := mr.reserves[id]; ok {
		return reserve, nil
	}
	return Reserve{}, ErrReserveNotFound
}

func (rep *MemoryRepository) GetByFilter(filter Reserve, ordered bool) (res map[uuid.UUID]Reserve, err error) {
	newrep := NewMemoryRepository(&rep.reserves, filter, ordered)
	return newrep.reserves, nil
}

func (rep *MemoryRepository) Add(r Reserve) (res Reserve, err error) {
	if rep.reserves == nil {
		rep.Lock()
		rep.reserves = make(map[uuid.UUID]Reserve)
		rep.Unlock()
	}

	if _, ok := rep.reserves[r.Id]; ok {
		err = fmt.Errorf("reserve already exists: %w", ErrFailedToAddReserve)
		return
	}
	rep.Lock()
	rep.reserves[r.Id] = r
	res = r
	rep.Unlock()
	return
}

func (mr *MemoryRepository) Update(memr Reserve) error {
	if _, ok := mr.reserves[memr.Id]; !ok {
		return fmt.Errorf("reserve does not exist: %w", ErrUpdateReserve)
	}
	mr.Lock()
	mr.reserves[memr.Id] = memr
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) AddPlayer(r Reserve, pl person.Person, count int) (Reserve, error) {
	r.Players[pl.Id] = Player{Person: pl, Count: count}
	return r, nil
}

func (mr *MemoryRepository) UpdatePlayer(r Reserve, pl person.Person, count int) (Reserve, error) {
	r.Players[pl.Id] = Player{Person: pl, Count: count}
	return r, nil
}
