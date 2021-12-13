package services

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"
)

func TestOrder_NewOrderService(t *testing.T) {

	os, err := NewOrderService(
		WithMemoryPersonRepository(),
		WithMemoryReserveRepository(),
	)

	if err != nil {
		t.Error(err)
	}

	p, err := person.NewPerson("Percy")
	if err != nil {
		t.Error(err)
	}

	err = os.persons.Add(p)
	if err != nil {
		t.Error(err)
	}

	duration, _ := time.ParseDuration("2h")
	err = os.CreateOrder(p.Id, time.Now(), time.Now().Add(duration))

	if err != nil {
		t.Error(err)
	}
}
