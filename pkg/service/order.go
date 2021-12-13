package services

import (
	"time"
	"volleybot/pkg/domain/person"
	permemory "volleybot/pkg/domain/person/memory"
	"volleybot/pkg/domain/reserve"
	resmemory "volleybot/pkg/domain/reserve/memory"

	"github.com/google/uuid"
)

type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	persons  person.PersonRepository
	reserves reserve.ReserveRepository
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
		os.persons = pr
		return nil
	}
}

func WithMemoryPersonRepository() OrderConfiguration {
	pr := permemory.New()
	return WithPersonRepository(pr)
}

func WithReserveRepository(rrep reserve.ReserveRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.reserves = rrep
		return nil
	}
}

func WithMemoryReserveRepository() OrderConfiguration {
	rrep := resmemory.New()
	return WithReserveRepository(rrep)
}

func (o *OrderService) CreateOrder(PersonId uuid.UUID, start time.Time, end time.Time) error {
	p, err := o.persons.Get(PersonId)
	if err != nil {
		return err
	}

	res, err := reserve.NewReserve(p, start, end)
	if err != nil {
		return err
	}
	o.reserves.Add(res)

	return nil
}
