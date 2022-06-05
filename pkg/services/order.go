package services

import (
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

type OrderConfiguration func(os *OrderService) error
type ReserveListResult struct {
	Reserves map[uuid.UUID]reserve.Reserve
	Err      error
}
type ReserveResult struct {
	Reserve reserve.Reserve
	Err     error
}

type OrderService struct {
	Persons   person.PersonRepository
	Reserves  reserve.ReserveRepository
	Locations location.LocationRepository
}

func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}
	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

func WithPersonRepository(pr person.PersonRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.Persons = pr
		return nil
	}
}

func WithMemoryPersonRepository() OrderConfiguration {
	pr := person.NewMemoryRepository()
	return WithPersonRepository(pr)
}

func WithPgPersonRepository(url string) OrderConfiguration {
	pr, _ := person.NewPgRepository(url)
	pr.UpdateDB()
	return WithPersonRepository(&pr)
}

func WithReserveRepository(rrep reserve.ReserveRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.Reserves = rrep
		return nil
	}
}

func WithMemoryReserveRepository() OrderConfiguration {
	rrep := reserve.NewMemoryRepository(nil, reserve.Reserve{}, false)
	return WithReserveRepository(&rrep)
}

func WithPgReserveRepository(url string) OrderConfiguration {
	rrep, _ := reserve.NewPgRepository(url)
	rrep.UpdateDB()
	return WithReserveRepository(&rrep)
}

func WithLocationRepository(rep location.LocationRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.Locations = rep
		return nil
	}
}

func WithPgLocationRepository(url string) OrderConfiguration {
	rep, _ := location.NewPgRepository(url)
	rep.UpdateDB()
	return WithLocationRepository(&rep)
}

func (serv *OrderService) CreateOrder(r reserve.Reserve, rchan chan ReserveResult) (res reserve.Reserve, err error) {
	res.Person, err = serv.Persons.Get(r.Person.Id)
	if err != nil {
		res.Person, err = person.NewPerson(r.Person.Firstname)
	}
	if err == nil {
		res, err = reserve.NewReserve(res.Person, r.StartTime, r.EndTime)
	}
	if err == nil {
		res.Location = r.Location
		serv.Reserves.Add(res)
	}
	if rchan != nil {
		rchan <- ReserveResult{Reserve: res, Err: err}
	}
	return
}

func (serv *OrderService) CancelOrder(r reserve.Reserve, rchan chan error) (err error) {
	r.Canceled = true
	err = serv.Reserves.Update(r)
	if rchan != nil {
		rchan <- err
	}
	return
}

func (serv *OrderService) List(filter reserve.Reserve, ordered bool, rchan chan ReserveListResult) (rlist map[uuid.UUID]reserve.Reserve, err error) {
	rlist, err = serv.Reserves.GetByFilter(filter, ordered)
	if rchan != nil {
		rchan <- ReserveListResult{Reserves: rlist, Err: err}
	}
	return
}
