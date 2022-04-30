package services

import (
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
	Persons  person.PersonRepository
	Reserves reserve.ReserveRepository
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
	rrep := reserve.NewMemoryRepository(nil, reserve.Reserve{})
	return WithReserveRepository(&rrep)
}

func WithPgReserveRepository(url string) OrderConfiguration {
	rrep, _ := reserve.NewPgRepository(url)
	rrep.UpdateDB()
	return WithReserveRepository(&rrep)
}

func (o *OrderService) CreateOrder(r reserve.Reserve, rchan chan ReserveResult) (result ReserveResult) {
	result.Reserve.Person, result.Err = o.Persons.Get(r.Person.Id)
	if result.Err != nil {
		result.Reserve.Person, result.Err = person.NewPerson(r.Person.Firstname)
	}
	if result.Err == nil {
		result.Reserve, result.Err = reserve.NewReserve(result.Reserve.Person, r.StartTime, r.EndTime)
	}
	if result.Err == nil {
		o.Reserves.Add(result.Reserve)
	}
	if rchan != nil {
		rchan <- result
	}
	return result
}

func (o *OrderService) List(filter reserve.Reserve, rchan chan ReserveListResult) ReserveListResult {
	rlist := ReserveListResult{}
	rlist.Reserves, rlist.Err = o.Reserves.GetByFilter(filter)
	if rchan != nil {
		rchan <- rlist
	}
	return rlist
}
