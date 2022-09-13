package volley

import (
	"fmt"
	"sync"
	"time"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	reserves []Volley
	sync.Mutex
}

func NewMemoryRepository(reserves *[]Volley, filter Volley, ordered bool) (mr MemoryRepository) {
	if reserves == nil {
		return MemoryRepository{
			reserves: []Volley{},
		}
	}

	mr.reserves = []Volley{}
	for _, v := range *reserves {
		if filter.Person.Id != uuid.Nil && v.Person.Id != filter.Person.Id {
			continue
		}

		if filter.StartTime != (time.Time{}) && filter.StartTime.After(v.EndTime) {
			continue
		}

		if filter.EndTime != (time.Time{}) && filter.EndTime.Before(v.StartTime) {
			continue
		}

		if filter.StartTime != (time.Time{}) && filter.EndTime != (time.Time{}) {
			fres := Volley{Reserve: reserve.Reserve{StartTime: filter.StartTime, EndTime: filter.EndTime}}
			if !fres.CheckConflicts(v.Reserve) {
				continue
			}
		}
		if ordered {
			if !v.Ordered() {
				continue
			}
		}

		mr.reserves = append(mr.reserves, v)
	}
	return
}

func (mr *MemoryRepository) Get(id uuid.UUID) (Volley, error) {
	for _, res := range mr.reserves {
		if res.Id == id {
			return res, nil
		}
	}
	return Volley{}, reserve.ErrReserveNotFound
}

func (rep *MemoryRepository) GetByFilter(filter Volley, ordered bool, sorted bool) (res []Volley, err error) {
	newrep := NewMemoryRepository(&rep.reserves, filter, ordered)
	return newrep.reserves, nil
}

func (rep *MemoryRepository) Add(r Volley) (res Volley, err error) {
	if rep.reserves == nil {
		rep.Lock()
		rep.reserves = []Volley{}
		rep.Unlock()
	}

	for _, rr := range rep.reserves {
		if rr.Id == r.Id {
			err = fmt.Errorf("reserve already exists: %w", reserve.ErrFailedToAddReserve)
			return
		}
	}
	rep.Lock()
	rep.reserves = append(rep.reserves, r)
	res = r
	rep.Unlock()
	return
}

func (mr *MemoryRepository) Update(memr Volley) error {
	for idx, res := range mr.reserves {
		if res.Id == memr.Id {
			mr.Lock()
			mr.reserves[idx] = memr
			mr.Unlock()
			return nil
		}
	}
	return fmt.Errorf("reserve does not exist: %w", reserve.ErrUpdateReserve)
}

func (mr *MemoryRepository) AddMember(r Volley, mb Member) (Volley, error) {
	for i, p := range r.Members {
		if p.Id == mb.Id {
			r.Members[i] = mb
			return r, nil
		}
	}
	r.Members = append(r.Members, mb)
	return r, nil
}

func (mr *MemoryRepository) UpdateMember(r Volley, mb Member) (Volley, error) {
	for i, p := range r.Members {
		if p.Id == mb.Id {
			r.Members[i] = mb
			return r, nil
		}
	}
	return r, reserve.ErrReservePlayerNotFound
}
